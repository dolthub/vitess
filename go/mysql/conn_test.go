/*
Copyright 2019 The Vitess Authors.

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
	crypto_rand "crypto/rand"
	"errors"
	"io"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func createSocketPair(t *testing.T) (net.Listener, *Conn, *Conn) {
	// Create a listener.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	addr := listener.Addr().String()
	listener.(*net.TCPListener).SetDeadline(time.Now().Add(10 * time.Second))

	// Dial a client, Accept a server.
	wg := sync.WaitGroup{}

	var clientConn net.Conn
	var clientErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		clientConn, clientErr = net.DialTimeout("tcp", addr, 10*time.Second)
	}()

	var serverConn net.Conn
	var serverErr error
	wg.Add(1)
	go func() {
		defer wg.Done()
		serverConn, serverErr = listener.Accept()
	}()

	wg.Wait()

	if clientErr != nil {
		t.Fatalf("Dial failed: %v", clientErr)
	}
	if serverErr != nil {
		t.Fatalf("Accept failed: %v", serverErr)
	}

	// Create a Conn on both sides.
	cConn := newConn(clientConn)
	sConn := newConn(serverConn)

	return listener, sConn, cConn
}

func useWritePacket(t *testing.T, cConn *Conn, data []byte) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatalf("%v", x)
		}
	}()
	if err := cConn.writePacket(data); err != nil {
		t.Fatalf("writePacket failed: %v", err)
	}
}

func useWriteEphemeralPacketBuffered(t *testing.T, cConn *Conn, data []byte) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatalf("%v", x)
		}
	}()
	cConn.startWriterBuffering()
	defer cConn.flush(context.Background())

	buf := cConn.startEphemeralPacket(len(data))
	copy(buf, data)
	if err := cConn.writeEphemeralPacket(); err != nil {
		t.Fatalf("writeEphemeralPacket(false) failed: %v", err)
	}
}

func useWriteEphemeralPacketDirect(t *testing.T, cConn *Conn, data []byte) {
	defer func() {
		if x := recover(); x != nil {
			t.Fatalf("%v", x)
		}
	}()

	buf := cConn.startEphemeralPacket(len(data))
	copy(buf, data)
	if err := cConn.writeEphemeralPacket(); err != nil {
		t.Fatalf("writeEphemeralPacket(true) failed: %v", err)
	}
}

func verifyPacketCommsSpecific(t *testing.T, cConn *Conn, data []byte,
	write func(t *testing.T, cConn *Conn, data []byte),
	read func(context.Context) ([]byte, error)) {
	// Have to do it in the background if it cannot be buffered.
	// Note we have to wait for it to finish at the end of the
	// test, as the write may write all the data to the socket,
	// and the flush may not be done after we're done with the read.
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		write(t, cConn, data)
		wg.Done()
	}()

	received, err := read(context.Background())
	if err != nil || !bytes.Equal(data, received) {
		t.Fatalf("ReadPacket failed: %v %v", received, err)
	}
	wg.Wait()
}

// Write a packet on one side, read it on the other, check it's
// correct.  We use all possible read and write methods.
func verifyPacketComms(t *testing.T, cConn, sConn *Conn, data []byte) {
	// All three writes, with ReadPacket.
	verifyPacketCommsSpecific(t, cConn, data, useWritePacket, sConn.ReadPacket)
	verifyPacketCommsSpecific(t, cConn, data, useWriteEphemeralPacketBuffered, sConn.ReadPacket)
	verifyPacketCommsSpecific(t, cConn, data, useWriteEphemeralPacketDirect, sConn.ReadPacket)

	// All three writes, with readEphemeralPacket.
	verifyPacketCommsSpecific(t, cConn, data, useWritePacket, sConn.readEphemeralPacket)
	sConn.recycleReadPacket()
	verifyPacketCommsSpecific(t, cConn, data, useWriteEphemeralPacketBuffered, sConn.readEphemeralPacket)
	sConn.recycleReadPacket()
	verifyPacketCommsSpecific(t, cConn, data, useWriteEphemeralPacketDirect, sConn.readEphemeralPacket)
	sConn.recycleReadPacket()

	// All three writes, with readEphemeralPacketDirect, if size allows it.
	if len(data) < MaxPacketSize {
		verifyPacketCommsSpecific(t, cConn, data, useWritePacket, sConn.readEphemeralPacketDirect)
		sConn.recycleReadPacket()
		verifyPacketCommsSpecific(t, cConn, data, useWriteEphemeralPacketBuffered, sConn.readEphemeralPacketDirect)
		sConn.recycleReadPacket()
		verifyPacketCommsSpecific(t, cConn, data, useWriteEphemeralPacketDirect, sConn.readEphemeralPacketDirect)
		sConn.recycleReadPacket()
	}
}

func TestPackets(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	// Verify all packets go through correctly.
	// Small one.
	data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	verifyPacketComms(t, cConn, sConn, data)

	// 0 length packet
	data = []byte{}
	verifyPacketComms(t, cConn, sConn, data)

	// Under the limit, still one packet.
	data = make([]byte, MaxPacketSize-1)
	data[0] = 0xab
	data[MaxPacketSize-2] = 0xef
	verifyPacketComms(t, cConn, sConn, data)

	// Exactly the limit, two packets.
	data = make([]byte, MaxPacketSize)
	data[0] = 0xab
	data[MaxPacketSize-1] = 0xef
	verifyPacketComms(t, cConn, sConn, data)

	// Over the limit, two packets.
	data = make([]byte, MaxPacketSize+1000)
	data[0] = 0xab
	data[MaxPacketSize+999] = 0xef
	verifyPacketComms(t, cConn, sConn, data)
}

func TestBasicPackets(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	// Write OK packet, read it, compare.
	if err := sConn.writeOKPacket(12, 34, 56, 78); err != nil {
		t.Fatalf("writeOKPacket failed: %v", err)
	}
	data, err := cConn.ReadPacket(context.Background())
	if err != nil || len(data) == 0 || data[0] != OKPacket {
		t.Fatalf("cConn.ReadPacket - OKPacket failed: %v %v", data, err)
	}
	affectedRows, lastInsertID, statusFlags, warnings, err := parseOKPacket(data)
	if err != nil || affectedRows != 12 || lastInsertID != 34 || statusFlags != 56 || warnings != 78 {
		t.Errorf("parseOKPacket returned unexpected data: %v %v %v %v %v", affectedRows, lastInsertID, statusFlags, warnings, err)
	}

	// Write OK packet with EOF header, read it, compare.
	if err := sConn.writeOKPacketWithEOFHeader(12, 34, 56, 78); err != nil {
		t.Fatalf("writeOKPacketWithEOFHeader failed: %v", err)
	}
	data, err = cConn.ReadPacket(context.Background())
	if err != nil || len(data) == 0 || !isEOFPacket(data) {
		t.Fatalf("cConn.ReadPacket - OKPacket with EOF header failed: %v %v", data, err)
	}
	affectedRows, lastInsertID, statusFlags, warnings, err = parseOKPacket(data)
	if err != nil || affectedRows != 12 || lastInsertID != 34 || statusFlags != 56 || warnings != 78 {
		t.Errorf("parseOKPacket returned unexpected data: %v %v %v %v %v", affectedRows, lastInsertID, statusFlags, warnings, err)
	}

	// Write error packet, read it, compare.
	if err := sConn.writeErrorPacket(ERAccessDeniedError, SSAccessDeniedError, "access denied: %v", "reason"); err != nil {
		t.Fatalf("writeErrorPacket failed: %v", err)
	}
	data, err = cConn.ReadPacket(context.Background())
	if err != nil || len(data) == 0 || data[0] != ErrPacket {
		t.Fatalf("cConn.ReadPacket - ErrorPacket failed: %v %v", data, err)
	}
	err = ParseErrorPacket(data)
	if !reflect.DeepEqual(err, NewSQLError(ERAccessDeniedError, SSAccessDeniedError, "access denied: reason")) {
		t.Errorf("ParseErrorPacket returned unexpected data: %v", err)
	}

	// Write error packet from error, read it, compare.
	if err := sConn.writeErrorPacketFromError(NewSQLError(ERAccessDeniedError, SSAccessDeniedError, "access denied")); err != nil {
		t.Fatalf("writeErrorPacketFromError failed: %v", err)
	}
	data, err = cConn.ReadPacket(context.Background())
	if err != nil || len(data) == 0 || data[0] != ErrPacket {
		t.Fatalf("cConn.ReadPacket - ErrorPacket failed: %v %v", data, err)
	}
	err = ParseErrorPacket(data)
	if !reflect.DeepEqual(err, NewSQLError(ERAccessDeniedError, SSAccessDeniedError, "access denied")) {
		t.Errorf("ParseErrorPacket returned unexpected data: %v", err)
	}

	// Write EOF packet, read it, compare first byte. Payload is always ignored.
	if err := sConn.writeEOFPacket(0x8912, 0xabba); err != nil {
		t.Fatalf("writeEOFPacket failed: %v", err)
	}
	data, err = cConn.ReadPacket(context.Background())
	if err != nil || len(data) == 0 || !isEOFPacket(data) {
		t.Fatalf("cConn.ReadPacket - EOFPacket failed: %v %v", data, err)
	}
}

// Mostly a sanity check.
func TestEOFOrLengthEncodedIntFuzz(t *testing.T) {
	for i := 0; i < 100; i++ {
		bytes := make([]byte, rand.Intn(16)+1)
		_, err := crypto_rand.Read(bytes)
		if err != nil {
			t.Fatalf("error doing rand.Read")
		}
		bytes[0] = 0xfe

		_, _, isInt := readLenEncInt(bytes, 0)
		isEOF := isEOFPacket(bytes)
		if (isInt && isEOF) || (!isInt && !isEOF) {
			t.Fatalf("0xfe bytestring is EOF xor Int. Bytes %v", bytes)
		}
	}
}

// waitForActivity runs sConn.WaitForClientActivity in a goroutine and returns
// its error, failing the test if it doesn't return within a generous timeout.
func waitForActivity(t *testing.T, sConn *Conn, ctx context.Context) error {
	t.Helper()
	errCh := make(chan error, 1)
	go func() {
		errCh <- sConn.WaitForClientActivity(ctx)
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("WaitForClientActivity did not return in time")
		return nil
	}
}

// TestWaitForClientActivityCancel verifies that cancelling the context unblocks
// the watch with a nil error and leaves the connection usable for a subsequent
// read.
func TestWaitForClientActivityCancel(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer listener.Close()
	defer sConn.Close()
	defer cConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- sConn.WaitForClientActivity(ctx)
	}()

	// Give the watch time to block in Peek, then cancel it.
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("expected nil after cancel, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("WaitForClientActivity did not return after cancel")
	}

	// The connection must still be usable: the deadline we set to interrupt the
	// peek must have been cleared.
	if _, err := cConn.Conn.Write([]byte{0x42}); err != nil {
		t.Fatalf("client write failed: %v", err)
	}
	b, err := sConn.bufferedReader.ReadByte()
	if err != nil {
		t.Fatalf("server read after cancel failed: %v", err)
	}
	if b != 0x42 {
		t.Fatalf("expected 0x42, got 0x%x", b)
	}
}

// TestWaitForClientActivityClosed verifies that a client disconnect unblocks the
// watch with a non-nil error.
func TestWaitForClientActivityClosed(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer listener.Close()
	defer sConn.Close()

	// Close the client end while the server is watching.
	go func() {
		time.Sleep(100 * time.Millisecond)
		cConn.Close()
	}()

	err := waitForActivity(t, sConn, context.Background())
	if err == nil {
		t.Fatal("expected non-nil error when client closed, got nil")
	}
	if errors.Is(err, ErrClientWroteWhileBusy) {
		t.Fatalf("expected a disconnect error, got ErrClientWroteWhileBusy")
	}
	// On a clean close the peer-side read observes io.EOF.
	if !errors.Is(err, io.EOF) {
		t.Logf("client-close error was %v (not io.EOF); acceptable if it's a reset", err)
	}
}

// TestWaitForClientActivityUnexpectedData verifies that data sent by the client
// while a query is "executing" is reported as ErrClientWroteWhileBusy and is not
// consumed (it remains available for the next command read).
func TestWaitForClientActivityUnexpectedData(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer listener.Close()
	defer sConn.Close()
	defer cConn.Close()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cConn.Conn.Write([]byte{0x07, 0x08, 0x09})
	}()

	err := waitForActivity(t, sConn, context.Background())
	if !errors.Is(err, ErrClientWroteWhileBusy) {
		t.Fatalf("expected ErrClientWroteWhileBusy, got %v", err)
	}

	// The bytes must not have been consumed by the watch.
	want := []byte{0x07, 0x08, 0x09}
	for i, w := range want {
		b, err := sConn.bufferedReader.ReadByte()
		if err != nil {
			t.Fatalf("reading byte %d after watch failed: %v", i, err)
		}
		if b != w {
			t.Fatalf("byte %d: expected 0x%x, got 0x%x", i, w, b)
		}
	}
}
