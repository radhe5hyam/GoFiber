package test

import (
	"net"
	"strings"
	"testing"

	"github.com/radhe5hyam/GoFiber/http"
)

type mockRequestConn struct {
	net.Conn
	reader *strings.Reader
}

func (m *mockRequestConn) Read(b []byte) (int, error) {
	return m.reader.Read(b)
}

func TestParseRawHTTPRequest(t *testing.T) {
	raw := "POST /submit?foo=bar HTTP/1.1\r\n" +
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
	if req.Headers["content-type"] != "application/x-www-form-urlencoded" {
		t.Errorf("Header mismatch: %+v", req.Headers)
	}
	if req.Query["foo"] != "bar" {
		t.Errorf("Expected query param foo=bar, got %v", req.Query)
	}
}

func TestParseHTTPRequest_NoBody(t *testing.T) {
	raw := "GET /hello HTTP/1.1\r\n" +
		"Host: localhost\r\n\r\n"

	mock := &mockRequestConn{reader: strings.NewReader(raw)}

	req, err := http.ParseHTTPRequest(mock)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("Expected GET, got %s", req.Method)
	}
	if req.Path != "/hello" {
		t.Errorf("Expected /hello, got %s", req.Path)
	}
	if len(req.Body) != 0 {
		t.Errorf("Expected empty body, got '%s'", string(req.Body))
	}
}

func TestParseHTTPRequest_MultipleHeaders(t *testing.T) {
	raw := "GET /test HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"X-Custom: value1\r\n" +
		"X-Another: value2\r\n\r\n"

	mock := &mockRequestConn{reader: strings.NewReader(raw)}

	req, err := http.ParseHTTPRequest(mock)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if req.Headers["x-custom"] != "value1" {
		t.Errorf("Expected x-custom header to be value1, got %s", req.Headers["x-custom"])
	}
	if req.Headers["x-another"] != "value2" {
		t.Errorf("Expected x-another header to be value2, got %s", req.Headers["x-another"])
	}
}

func TestParseHTTPRequest_QueryParams(t *testing.T) {
	raw := "GET /search?q=go+fiber&lang=en HTTP/1.1\r\n" +
		"Host: localhost\r\n\r\n"

	mock := &mockRequestConn{reader: strings.NewReader(raw)}

	req, err := http.ParseHTTPRequest(mock)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if req.Query["q"] != "go fiber" {
		t.Errorf("Expected query param q=go fiber, got %v", req.Query["q"])
	}
	if req.Query["lang"] != "en" {
		t.Errorf("Expected query param lang=en, got %v", req.Query["lang"])
	}
}

func TestParseHTTPRequest_InvalidRequestLine(t *testing.T) {
	raw := "INVALIDREQUEST\r\nHost: localhost\r\n\r\n"

	mock := &mockRequestConn{reader: strings.NewReader(raw)}

	_, err := http.ParseHTTPRequest(mock)
	if err == nil {
		t.Errorf("Expected error for invalid request line, got nil")
	}
}

func TestParseHTTPRequest_InvalidContentLength(t *testing.T) {
	raw := "POST / HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Content-Length: abc\r\n\r\n" +
		"body"

	mock := &mockRequestConn{reader: strings.NewReader(raw)}

	req, err := http.ParseHTTPRequest(mock)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	// Should not panic, body should be empty since content-length is invalid
	if len(req.Body) != 0 {
		t.Errorf("Expected empty body for invalid content-length, got '%s'", string(req.Body))
	}
}

func TestParseHTTPRequest_EmptyHeaders(t *testing.T) {
	raw := "GET /empty HTTP/1.1\r\n\r\n"

	mock := &mockRequestConn{reader: strings.NewReader(raw)}

	req, err := http.ParseHTTPRequest(mock)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if req.Path != "/empty" {
		t.Errorf("Expected /empty, got %s", req.Path)
	}
	if len(req.Headers) != 0 {
		t.Errorf("Expected no headers, got %+v", req.Headers)
	}
}
