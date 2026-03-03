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
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/dolthub/vitess/go/vt/tlstest"
	"github.com/dolthub/vitess/go/vt/vttls"
)

// This file tests the handshake scenarios between our client and our server.

func TestClearTextClientAuth(t *testing.T) {
	th := &testHandler{}

	authServer := NewAuthServerStaticWithAuthMethodDescription("", "", 0, MysqlClearPassword)
	authServer.entries["user1"] = []*AuthServerStaticEntry{
		{Password: "password1"},
	}

	// Create the listener.
	l, err := NewListener("tcp", ":0", authServer, th, 0, 0)
	if err != nil {
		t.Fatalf("NewListener failed: %v", err)
	}
	defer l.Close()
	host := l.Addr().(*net.TCPAddr).IP.String()
	port := l.Addr().(*net.TCPAddr).Port
	go func() {
		l.Accept()
	}()

	// Setup the right parameters.
	params := &ConnParams{
		Host:  host,
		Port:  port,
		Uname: "user1",
		Pass:  "password1",
	}

	// Connection should fail, as server requires SSL for clear text auth.
	ctx := context.Background()
	_, err = Connect(ctx, params)
	if err == nil {
		t.Fatalf("expected connection error")
	}
	sqlErr, ok := err.(*SQLError)
	if !ok {
		t.Fatalf("expected *SQLError, got: %T (%v)", err, err)
	}
	if sqlErr.Number() != ERAccessDeniedError {
		t.Fatalf("unexpected mysql error code: %d", sqlErr.Number())
	}
	if sqlErr.SQLState() != SSAccessDeniedError {
		t.Fatalf("unexpected sqlstate: %s", sqlErr.SQLState())
	}

	// Change server side to allow clear text without auth.
	l.AllowClearTextWithoutTLS.Set(true)
	conn, err := Connect(ctx, params)
	if err != nil {
		t.Fatalf("unexpected connection error: %v", err)
	}
	defer conn.Close()

	// Run a 'select rows' command with results.
	result, err := conn.ExecuteFetch("select rows", 10000, true)
	if err != nil {
		t.Fatalf("ExecuteFetch failed: %v", err)
	}
	if !reflect.DeepEqual(result, selectRowsResult) {
		t.Errorf("Got wrong result from ExecuteFetch(select rows): %v", result)
	}

	// Send a ComQuit to avoid the error message on the server side.
	conn.writeComQuit()
}

// TestSSLConnection creates a server with TLS support, a client that
// also has SSL support, and connects them.
func TestSSLConnection(t *testing.T) {
	th := &testHandler{}

	authServer := NewAuthServerStaticWithAuthMethodDescription("", "", 0, MysqlClearPassword)
	authServer.entries["user1"] = []*AuthServerStaticEntry{
		{Password: "password1"},
	}

	// Create the listener, so we can get its host.
	l, err := NewListener("tcp", ":0", authServer, th, 0, 0)
	if err != nil {
		t.Fatalf("NewListener failed: %v", err)
	}
	defer l.Close()
	host := l.Addr().(*net.TCPAddr).IP.String()
	port := l.Addr().(*net.TCPAddr).Port

	// Create the certs.
	root, err := ioutil.TempDir("", "TestSSLConnection")
	if err != nil {
		t.Fatalf("TempDir failed: %v", err)
	}
	defer os.RemoveAll(root)
	tlstest.CreateCA(root)
	tlstest.CreateSignedCert(root, tlstest.CA, "01", "server", "server.example.com")
	tlstest.CreateSignedCert(root, tlstest.CA, "02", "client", "Client Cert")

	// Create the server with TLS config.
	serverConfig, err := vttls.ServerConfig(
		path.Join(root, "server-cert.pem"),
		path.Join(root, "server-key.pem"),
		path.Join(root, "ca-cert.pem"),
		"",
		"",
		tls.VersionTLS12)
	if err != nil {
		t.Fatalf("TLSServerConfig failed: %v", err)
	}
	l.TLSConfig = serverConfig
	go func() {
		l.Accept()
	}()

	// Setup the right parameters.
	params := &ConnParams{
		Host:  host,
		Port:  port,
		Uname: "user1",
		Pass:  "password1",
		// SSL flags.
		Flags:      CapabilityClientSSL,
		SslCa:      path.Join(root, "ca-cert.pem"),
		SslCert:    path.Join(root, "client-cert.pem"),
		SslKey:     path.Join(root, "client-key.pem"),
		ServerName: "server.example.com",
	}

	t.Run("Basics", func(t *testing.T) {
		testSSLConnectionBasics(t, params)
	})

	// Make sure clear text auth works over SSL.
	t.Run("ClearText", func(t *testing.T) {
		testSSLConnectionClearText(t, params)
	})
}

func testSSLConnectionClearText(t *testing.T, params *ConnParams) {
	// Create a client connection, connect.
	ctx := context.Background()
	conn, err := Connect(ctx, params)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer conn.Close()
	if conn.User != "user1" {
		t.Errorf("Invalid conn.User, got %v was expecting user1", conn.User)
	}

	// Make sure this went through SSL.
	result, err := conn.ExecuteFetch("ssl echo", 10000, true)
	if err != nil {
		t.Fatalf("ExecuteFetch failed: %v", err)
	}
	if result.Rows[0][0].ToString() != "ON" {
		t.Errorf("Got wrong result from ExecuteFetch(ssl echo): %v", result)
	}

	// Send a ComQuit to avoid the error message on the server side.
	conn.writeComQuit()
}

func testSSLConnectionBasics(t *testing.T, params *ConnParams) {
	// Create a client connection, connect.
	ctx := context.Background()
	conn, err := Connect(ctx, params)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer conn.Close()
	if conn.User != "user1" {
		t.Errorf("Invalid conn.User, got %v was expecting user1", conn.User)
	}

	// Run a 'select rows' command with results.
	result, err := conn.ExecuteFetch("select rows", 10000, true)
	if err != nil {
		t.Fatalf("ExecuteFetch failed: %v", err)
	}
	if !reflect.DeepEqual(result, selectRowsResult) {
		t.Errorf("Got wrong result from ExecuteFetch(select rows): %v", result)
	}

	// Make sure this went through SSL.
	result, err = conn.ExecuteFetch("ssl echo", 10000, true)
	if err != nil {
		t.Fatalf("ExecuteFetch failed: %v", err)
	}
	if result.Rows[0][0].ToString() != "ON" {
		t.Errorf("Got wrong result from ExecuteFetch(ssl echo): %v", result)
	}

	// Send a ComQuit to avoid the error message on the server side.
	conn.writeComQuit()
}

// rejectingUserValidator is a [UserValidator] that rejects every user.
// Embedding it in an auth server drives the negotiated-method == nil branch.
type rejectingUserValidator struct{}

func (rejectingUserValidator) HandleUser(string, net.Addr) bool { return false }

// rejectingHashStorage is a minimal [HashStorage] that always returns access denied.
type rejectingHashStorage struct{}

func (rejectingHashStorage) UserEntryWithHash(_ *Conn, _ []byte, user string, _ []byte, _ net.Addr) (Getter, error) {
	return nil, NewSQLError(ERAccessDeniedError, SSAccessDeniedError, "Access denied for user '%v'", user)
}

// rejectAllAuthServer is an [AuthServer] whose sole auth method always returns false
// from HandleUser, so no method is ever negotiated for any user.
type rejectAllAuthServer struct{}

func (s *rejectAllAuthServer) AuthMethods() []AuthMethod {
	return []AuthMethod{NewMysqlNativeAuthMethod(rejectingHashStorage{}, rejectingUserValidator{})}
}

func (s *rejectAllAuthServer) DefaultAuthMethodDescription() AuthMethodDescription {
	return MysqlNativePassword
}

// TestNoAuthMethodsReturnsAccessDenied exercises the server branch where
// negotiatedAuthMethod is nil after both the initial negotiation and the
// fallback scan of all methods fail.
//
// Before the fix, the server sent CRServerHandshakeErr (2012), a client-side
// error code that strict clients (e.g. MariaDB 11.8+) interpret as a malformed
// packet (CR_MALFORMED_PACKET/2027). The server must instead emit a
// spec-compliant ERR_Packet with ERAccessDeniedError (1045) and
// SSAccessDeniedError ("28000").
func TestNoAuthMethodsReturnsAccessDenied(t *testing.T) {
	th := &testHandler{}
	authServer := &rejectAllAuthServer{}
	l, err := NewListener("tcp", ":0", authServer, th, 0, 0)
	if err != nil {
		t.Fatalf("NewListener failed: %v", err)
	}
	defer l.Close()
	host := l.Addr().(*net.TCPAddr).IP.String()
	port := l.Addr().(*net.TCPAddr).Port
	go l.Accept()

	params := &ConnParams{
		Host:  host,
		Port:  port,
		Uname: "nonexistent",
		Pass:  "anything",
	}
	ctx := context.Background()
	_, err = Connect(ctx, params)
	if err == nil {
		t.Fatal("expected connection error, got none")
	}
	sqlErr, ok := err.(*SQLError)
	if !ok {
		t.Fatalf("expected *SQLError, got: %T (%v)", err, err)
	}
	if sqlErr.Number() != ERAccessDeniedError {
		t.Fatalf("expected error code %d (ERAccessDeniedError), got %d", ERAccessDeniedError, sqlErr.Number())
	}
	if sqlErr.SQLState() != SSAccessDeniedError {
		t.Fatalf("expected SQL state %q (SSAccessDeniedError), got %q", SSAccessDeniedError, sqlErr.SQLState())
	}
}
