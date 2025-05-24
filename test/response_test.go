package test

import (
	"bufio"
	"bytes"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/radhe5hyam/GoFiber/http"
)

type mockResponseConn struct {
	reader *strings.Reader
	buffer *bytes.Buffer
}
func (m *mockResponseConn) LocalAddr() net.Addr              { return nil }
func (m *mockResponseConn) RemoteAddr() net.Addr             { return nil }
func (m *mockResponseConn) SetDeadline(t time.Time) error    { return nil }
func (m *mockResponseConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockResponseConn) SetWriteDeadline(t time.Time) error { return nil }

func (m *mockResponseConn) Read(b []byte) (int, error) {
	return m.reader.Read(b)
}
func (m *mockResponseConn) Write(b []byte) (int, error) {
	return m.buffer.Write(b)
}
func (m *mockResponseConn) Close() error { return nil }

func TestHTTPResponseWriter(t *testing.T) {
	buffer := &bytes.Buffer{}
	conn := &mockResponseConn{buffer: buffer}

	writer := http.NewResponseWriter(conn)
	writer.Header("Content-Type", "text/plain")
	writer.WriteHeader(200)
	writer.Write([]byte("Hello, World!"))
	err := writer.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	scanner := bufio.NewScanner(buffer)
	scanner.Scan()
	statusLine := scanner.Text()
	if statusLine != "HTTP/1.1 200 OK" {
		t.Errorf("Status line mismatch: %s", statusLine)
	}

	foundContentType := false
	foundContentLength := false
	for scanner.Scan() {
		line := scanner.Text()
		t.Log("Header line:", line)
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Type:") {
			foundContentType = true
		}
		if strings.HasPrefix(line, "Content-Length:") {
			t.Log("Content-Length header found:", line)
			foundContentLength = true
		}
	}

	if !foundContentType || !foundContentLength {
		t.Errorf("Missing headers: Content-Type=%v, Content-Length=%v", foundContentType, foundContentLength)
	}
}
