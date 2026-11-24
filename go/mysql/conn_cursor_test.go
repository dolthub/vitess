/*
Copyright 2026 Dolthub, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package mysql

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/dolthub/vitess/go/sqltypes"
	querypb "github.com/dolthub/vitess/go/vt/proto/query"
)

// TestWaitForClientActivityWhileCursorOpen checks that opening a cursor preempts
// a watcher parked in Peek, that the watcher neither consumes nor reports the
// client's commands while the cursor is open, and that it stays out of the read
// side until its query ends even after the cursor closes.
func TestWaitForClientActivityWhileCursorOpen(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer listener.Close()
	defer sConn.Close()
	defer cConn.Close()

	watchCtx, cancelWatch := context.WithCancel(context.Background())
	defer cancelWatch()
	watchErr := make(chan error, 1)
	go func() { watchErr <- sConn.WaitForClientActivity(watchCtx) }()
	time.Sleep(100 * time.Millisecond)

	sConn.readLease.markCursorOpen()

	want := []byte{ComStmtFetch, 0x01, 0x00, 0x00, 0x00}
	if _, err := cConn.Conn.Write(want); err != nil {
		t.Fatalf("client write failed: %v", err)
	}
	select {
	case err := <-watchErr:
		t.Fatalf("watcher returned while the cursor was open: %v", err)
	case <-time.After(200 * time.Millisecond):
	}
	got := make([]byte, len(want))
	if err := sConn.Conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadFull(sConn.bufferedReader, got); err != nil {
		t.Fatalf("command read failed: %v", err)
	}
	if err := sConn.Conn.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("command bytes corrupted: got %v, want %v", got, want)
	}

	sConn.readLease.markCursorClosed()
	if _, err := cConn.Conn.Write([]byte{0x01}); err != nil {
		t.Fatalf("client write failed: %v", err)
	}
	select {
	case err := <-watchErr:
		t.Fatalf("watcher returned after the cursor closed but before its query ended: %v", err)
	case <-time.After(200 * time.Millisecond):
	}
	cancelWatch()
	select {
	case err := <-watchErr:
		if err != nil {
			t.Fatalf("expected nil once the query ended, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("watcher did not return once its query ended")
	}

	// A watcher for a query started after the cursor closed watches as usual.
	go func() { watchErr <- sConn.WaitForClientActivity(context.Background()) }()
	select {
	case err := <-watchErr:
		if !errors.Is(err, ErrClientWroteWhileBusy) {
			t.Fatalf("expected ErrClientWroteWhileBusy, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("watcher did not report the buffered client byte")
	}
}

// cursorTestHandler streams a fixed number of rows in batches from
// ComStmtExecute and, as go-mysql-server does for every query, watches the
// client connection for as long as the statement executes.
type cursorTestHandler struct {
	*testHandler
	rows, batch int
	watchErr    chan error
}

func (h *cursorTestHandler) ComStmtExecute(ctx context.Context, c *Conn, prepare *PrepareData, callback func(*sqltypes.Result) error) error {
	watchCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() { h.watchErr <- c.WaitForClientActivity(watchCtx) }()

	fields := []*querypb.Field{{Name: "id", Type: querypb.Type_INT32}}
	for start := 1; start <= h.rows; start += h.batch {
		res := &sqltypes.Result{Fields: fields}
		for id := start; id < start+h.batch && id <= h.rows; id++ {
			res.Rows = append(res.Rows, []sqltypes.Value{
				sqltypes.MakeTrusted(querypb.Type_INT32, []byte(strconv.Itoa(id))),
			})
		}
		if start+h.batch > h.rows {
			// go-mysql-server updates the session flags here, between the last row
			// and the final callback, while the cursor is still open.
			c.StatusFlags |= uint16(ServerStatusAutocommit)
		}
		if err := callback(res); err != nil {
			return err
		}
	}
	return nil
}

// TestServerSideCursorFetch drives a read-only cursor through handleNextCommand
// with the handler watching the connection. The watcher must not fire on the
// client's fetches, each fetch that leaves the cursor open must end in an EOF
// carrying SERVER_STATUS_CURSOR_EXISTS, and the fetch that exhausts it must end
// in an EOF carrying SERVER_STATUS_LAST_ROW_SENT alone, as MySQL does.
func TestServerSideCursorFetch(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer listener.Close()
	defer sConn.Close()
	defer cConn.Close()

	handler := &cursorTestHandler{
		testHandler: &testHandler{},
		rows:        5,
		batch:       2,
		watchErr:    make(chan error, 1),
	}
	const stmtID = 7
	sConn.PrepareData = map[uint32]*PrepareData{
		stmtID: {StatementID: stmtID, PrepareStmt: "select id from t"},
	}

	ctx := context.Background()
	serverErr := make(chan error, 1)
	serve := func() { serverErr <- sConn.handleNextCommand(ctx, handler) }
	send := func(packet []byte) {
		t.Helper()
		cConn.sequence = 0
		if err := cConn.writePacket(packet); err != nil {
			t.Fatalf("client write failed: %v", err)
		}
		if err := cConn.flush(ctx); err != nil {
			t.Fatalf("client flush failed: %v", err)
		}
	}
	// readResponse reads packets up to the EOF that ends a response and returns
	// how many packets preceded it along with the EOF's status flags.
	readResponse := func() (int, serverStatus) {
		t.Helper()
		n := 0
		for {
			if err := cConn.Conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
				t.Fatal(err)
			}
			data, err := cConn.ReadPacket(ctx)
			if err != nil {
				t.Fatalf("client read failed: %v", err)
			}
			if isEOFPacket(data) {
				_, status, err := parseEOFPacket(data)
				if err != nil {
					t.Fatalf("parsing EOF failed: %v", err)
				}
				return n, status
			}
			if data[0] == ErrPacket {
				t.Fatalf("server returned an error packet: %v", data)
			}
			n++
		}
	}
	cursorOpen := func() bool {
		sConn.readLease.mu.Lock()
		defer sConn.readLease.mu.Unlock()
		return sConn.readLease.cursorOpen
	}

	// COM_STMT_EXECUTE with CURSOR_TYPE_READ_ONLY, iteration count 1, no params.
	go serve()
	send([]byte{ComStmtExecute, stmtID, 0, 0, 0, ReadOnly, 1, 0, 0, 0})
	packets, status := readResponse()
	if err := <-serverErr; err != nil {
		t.Fatalf("handleNextCommand(execute) failed: %v", err)
	}
	if packets != 2 { // column count + one column definition, no rows
		t.Fatalf("execute: expected 2 packets before EOF, got %d", packets)
	}
	if status&ServerCursorExists == 0 {
		t.Fatalf("execute EOF must carry SERVER_STATUS_CURSOR_EXISTS, got %#x", status)
	}
	if !cursorOpen() {
		t.Fatal("the read lease must mark the cursor open")
	}

	// Fetch two rows at a time: five rows means two fetches leave the cursor open
	// and the third exhausts it.
	fetches := []struct {
		rows   int
		status serverStatus
	}{
		{2, ServerCursorExists},
		{2, ServerCursorExists},
		{1, ServerCursorLastRowSent},
	}
	for i, want := range fetches {
		// Give an unguarded watcher time to park in Peek, as go-mysql-server's does.
		time.Sleep(100 * time.Millisecond)
		go serve()
		send([]byte{ComStmtFetch, stmtID, 0, 0, 0, 2, 0, 0, 0})
		rows, status := readResponse()
		if err := <-serverErr; err != nil {
			t.Fatalf("handleNextCommand(fetch %d) failed: %v", i+1, err)
		}
		if rows != want.rows {
			t.Fatalf("fetch %d: got %d rows, want %d", i+1, rows, want.rows)
		}
		if got := status & (ServerCursorExists | ServerCursorLastRowSent); got != want.status {
			t.Fatalf("fetch %d: cursor status flags %#x, want %#x", i+1, got, want.status)
		}
	}
	if sConn.cs != nil {
		t.Fatal("cursor state must be cleared once the last row was sent")
	}
	if cursorOpen() {
		t.Fatal("the read lease must mark the cursor closed once the last row was sent")
	}
	if sConn.StatusFlags&uint16(ServerStatusAutocommit) == 0 {
		t.Fatal("the handler's status flag update was lost")
	}
	if got := sConn.StatusFlags & uint16(ServerCursorExists|ServerCursorLastRowSent); got != 0 {
		t.Fatalf("cursor status flags leaked into the connection's flags: %#x", got)
	}

	// The handler's watcher returned because the query finished, not because
	// it mistook a fetch for a departed client.
	select {
	case err := <-handler.watchErr:
		if err != nil {
			t.Fatalf("watcher fired on the cursor's fetches: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("watcher did not return after the cursor query finished")
	}
}
