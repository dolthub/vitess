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

package netutil

import (
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func createSocketPair(t *testing.T) (net.Listener, net.Conn, net.Conn) {
	// Create a listener.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	addr := listener.Addr().String()

	// Dial a client, Accept a server.
	wg := sync.WaitGroup{}

	var clientConn net.Conn
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		clientConn, err = net.Dial("tcp", addr)
		if err != nil {
			t.Errorf("Dial failed: %v", err)
		}
	}()

	var serverConn net.Conn
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		serverConn, err = listener.Accept()
		if err != nil {
			t.Errorf("Accept failed: %v", err)
		}
	}()

	wg.Wait()

	return listener, serverConn, clientConn
}

func TestReadTimeout(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	cConnWithTimeout := NewConnWithTimeouts(cConn, 1*time.Millisecond, 1*time.Millisecond)

	c := make(chan error, 1)
	go func() {
		_, err := cConnWithTimeout.Read(make([]byte, 10))
		c <- err
	}()

	select {
	case err := <-c:
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		if !strings.HasSuffix(err.Error(), "i/o timeout") {
			t.Errorf("Expected error timeout, got %s", err)
		}
	case <-time.After(10 * time.Second):
		t.Errorf("Timeout did not happen")
	}
}

func TestWriteTimeout(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	sConnWithTimeout := NewConnWithTimeouts(sConn, 1*time.Millisecond, 1*time.Millisecond)

	c := make(chan error, 1)
	go func() {
		// The timeout will trigger when the buffer is full, so to test this we need to write multiple times.
		for {
			_, err := sConnWithTimeout.Write([]byte("payload"))
			if err != nil {
				c <- err
				return
			}
		}
	}()

	select {
	case err := <-c:
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		if !strings.HasSuffix(err.Error(), "i/o timeout") {
			t.Errorf("Expected error timeout, got %s", err)
		}
	case <-time.After(10 * time.Second):
		t.Errorf("Timeout did not happen")
	}
}

func TestNoTimeouts(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	cConnWithTimeout := NewConnWithTimeouts(cConn, 0, 24*time.Hour)

	c := make(chan error, 1)
	go func() {
		_, err := cConnWithTimeout.Read(make([]byte, 10))
		c <- err
	}()

	select {
	case <-c:
		t.Fatalf("Connection timeout, without a timeout")
	case <-time.After(100 * time.Millisecond):
		// NOOP
	}

	c2 := make(chan error, 1)
	sConnWithTimeout := NewConnWithTimeouts(sConn, 24*time.Hour, 0)
	go func() {
		// This should not fail as there is not timeout on write.
		for {
			_, err := sConnWithTimeout.Write([]byte("payload"))
			if err != nil {
				c2 <- err
				return
			}
		}
	}()
	select {
	case <-c2:
		t.Fatalf("Connection timeout, without a timeout")
	case <-time.After(100 * time.Millisecond):
		// NOOP
	}
}

// TestSingleShotOverrideResumesManagedTimeout verifies the single-shot override
// contract: SetReadDeadline with an explicit (non-zero) deadline applies to the
// next Read only, which skips the managed re-arm; the managed timeout then
// resumes on the following Read. This guards that a transient deadline override
// (e.g. a watcher cancelling a blocked Read) does not permanently disable the
// connection's managed timeout, so net_read_timeout still reaps afterward.
func TestSingleShotOverrideResumesManagedTimeout(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	// Short managed read timeout; the peer never sends, so a managed Read times
	// out quickly while an overridden Read with a far-future deadline does not.
	const managed = 50 * time.Millisecond
	s := NewConnWithTimeouts(sConn, managed, 0)

	// Single-shot override: a far-future explicit deadline for the next Read.
	if err := s.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline override failed: %v", err)
	}
	c := make(chan error, 1)
	go func() {
		_, err := s.Read(make([]byte, 10))
		c <- err
	}()

	// The override must suppress the managed re-arm: the Read stays blocked well
	// past the managed timeout.
	select {
	case err := <-c:
		t.Fatalf("overridden Read returned (err=%v); managed timeout was not skipped", err)
	case <-time.After(10 * managed):
		// Good: still blocked under the far-future override.
	}

	// Wake the blocked Read with a past deadline, then resume managed timeouts.
	if err := s.SetReadDeadline(time.Now().Add(-time.Hour)); err != nil {
		t.Fatalf("SetReadDeadline (wake) failed: %v", err)
	}
	select {
	case err := <-c:
		if err == nil {
			t.Fatalf("expected timeout error waking the overridden Read, got nil")
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("overridden Read did not wake after past deadline")
	}
	// Zero time resumes the managed timeout (drops the override).
	if err := s.SetReadDeadline(time.Time{}); err != nil {
		t.Fatalf("SetReadDeadline(zero) failed: %v", err)
	}

	// The next Read must re-arm the managed timeout again (reaping resumed).
	c2 := make(chan error, 1)
	go func() {
		_, err := s.Read(make([]byte, 10))
		c2 <- err
	}()
	select {
	case err := <-c2:
		if err == nil || !strings.HasSuffix(err.Error(), "i/o timeout") {
			t.Fatalf("expected managed i/o timeout after resume, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("managed timeout did not resume after a single-shot override")
	}
}
