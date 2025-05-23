package test

import (
	"net"
	"strings"
	"testing"

	"github.com/radhe5hyam/GoFiber/http"
)

type mockRequestConn  struct {
	net.Conn
	reader *strings.Reader
}

func (m *mockRequestConn) Read(b []byte) (int, error) {
	return m.reader.Read(b)
}

func TestParseRawHTTPRequest(t *testing.T) {
	raw := "POST /submit HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Content-Length: 11\r\n" +
		"Content-Type: application/x-www-form-urlencoded\r\n\r\n" +
		"hello=world"

	mock := &mockRequestConn{reader: strings.NewReader(raw)}

	req, err := http.ParseHTTPRequest(mock)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("Expected POST, got %s", req.Method)
	}
	if req.Path != "/submit" {
		t.Errorf("Expected /submit, got %s", req.Path)
	}
	if req.Headers["Content-Type"] != "application/x-www-form-urlencoded" {
		t.Errorf("Header mismatch: %+v", req.Headers)
	}
	if req.Body != "hello=world" {
		t.Errorf("Expected body 'hello=world', got '%s'", req.Body)
	}
}
