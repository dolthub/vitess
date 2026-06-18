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
	"sync"
	"time"
)

var _ net.Conn = (*ConnWithTimeouts)(nil)

// A ConnWithTimeouts is a wrapper to net.Conn that allows to set a read and write timeouts.
type ConnWithTimeouts struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
	// readOverride/writeOverride, when true, indicate that an explicit deadline
	// has been set via Set{,Read,Write}Deadline and must be honored for exactly
	// the next Read/Write instead of re-arming the managed timeout. They are
	// cleared when that next operation consumes them, after which the managed
	// timeout resumes.
	readOverride  bool
	writeOverride bool
	mu            sync.Mutex
}

// NewConnWithTimeouts wraps a net.Conn with managed read and write timeouts.
// It re-arms a fresh read or write deadline (now + the configured duration) at
// the start of every Read or Write, so each operation is bounded by that
// duration. A zero duration disables the corresponding managed timeout.
//
// Set{,Read,Write}Deadline with a non-zero time is a single-shot override: the
// explicit deadline is forwarded to the underlying connection and applies to the
// next Read/Write only (which skips the managed re-arm); the managed timeout then
// resumes on the following operation. Set{,Read,Write}Deadline with the zero time
// drops any pending override and resumes the managed timeout immediately. This
// lets a caller transiently interrupt or extend a single operation (e.g. cancel a
// blocked Read by setting a deadline in the past) without permanently disabling
// the connection's managed timeout.
//
// The returned *ConnWithTimeouts is not safe to copy.
func NewConnWithTimeouts(conn net.Conn, readTimeout time.Duration, writeTimeout time.Duration) *ConnWithTimeouts {
	return &ConnWithTimeouts{Conn: conn, readTimeout: readTimeout, writeTimeout: writeTimeout}
}

// Implementation of the Conn interface.

// Read sets a read deadline and delegates to conn.Read.
func (c *ConnWithTimeouts) Read(b []byte) (int, error) {
	c.mu.Lock()
	if c.readOverride {
		// An explicit deadline was set for this read; honor it once, then let
		// the managed timeout resume on the next read.
		c.readOverride = false
		c.mu.Unlock()
		return c.Conn.Read(b)
	}
	if c.readTimeout == 0 {
		c.mu.Unlock()
		return c.Conn.Read(b)
	}
	if err := c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
		c.mu.Unlock()
		return 0, err
	}
	c.mu.Unlock()
	return c.Conn.Read(b)
}

// Write sets a write deadline and delegates to conn.Write
func (c *ConnWithTimeouts) Write(b []byte) (int, error) {
	c.mu.Lock()
	if c.writeOverride {
		// An explicit deadline was set for this write; honor it once, then let
		// the managed timeout resume on the next write.
		c.writeOverride = false
		c.mu.Unlock()
		return c.Conn.Write(b)
	}
	if c.writeTimeout == 0 {
		c.mu.Unlock()
		return c.Conn.Write(b)
	}
	if err := c.Conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
		c.mu.Unlock()
		return 0, err
	}
	c.mu.Unlock()
	return c.Conn.Write(b)
}

// SetDeadline implements the Conn SetDeadline method. See the type doc for the
// single-shot override / zero-resumes-managed semantics.
func (c *ConnWithTimeouts) SetDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.readOverride = !t.IsZero()
	c.writeOverride = !t.IsZero()
	return c.Conn.SetDeadline(t)
}

// SetReadDeadline implements the Conn SetReadDeadline method. See the type doc
// for the single-shot override / zero-resumes-managed semantics.
func (c *ConnWithTimeouts) SetReadDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.readOverride = !t.IsZero()
	return c.Conn.SetReadDeadline(t)
}

// SetWriteDeadline implements the Conn SetWriteDeadline method. See the type doc
// for the single-shot override / zero-resumes-managed semantics.
func (c *ConnWithTimeouts) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.writeOverride = !t.IsZero()
	return c.Conn.SetWriteDeadline(t)
}
