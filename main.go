package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)


type connResponseWriter struct {
	conn           net.Conn
	header         http.Header
	statusCode     int
	headersWritten bool
}

func newConnResponseWriter(conn net.Conn) *connResponseWriter {
	return &connResponseWriter{
		conn:       conn,
		header:     make(http.Header),
		statusCode: http.StatusOK,
	}
}

func (crw *connResponseWriter) Header() http.Header {
	return crw.header
}

func (crw *connResponseWriter) Write(data []byte) (int, error) {
	if !crw.headersWritten {
		crw.WriteHeader(crw.statusCode)
	}
	return crw.conn.Write(data)
}

func (crw *connResponseWriter) WriteHeader(statusCode int) {
	if crw.headersWritten {
		return
	}
	crw.statusCode = statusCode

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", crw.statusCode, http.StatusText(crw.statusCode))
	crw.conn.Write([]byte(statusLine))

	if crw.header.Get("Content-Type") == "" {
		crw.header.Set("Content-Type", "text/plain; charset=utf-8")
	}

	for k, v := range crw.header {
		for _, val := range v {
			headerLine := fmt.Sprintf("%s: %s\r\n", k, val)
			crw.conn.Write([]byte(headerLine))
		}
	}
	crw.conn.Write([]byte("\r\n"))
	crw.headersWritten = true
}

// LoggerMiddleware is a simple example of a middleware.
func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("--> %s %s", r.Method, r.URL.Path)
		next(w, r) // Call the next handler in the chain
		log.Printf("<-- %s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
}

func main() {
	router := NewRouter()
	// Register global middleware
	router.Use(LoggerMiddleware)
	SetupExampleRoutes(router)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error setting up listener:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server started on port 8080, now using the custom router!")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleRequest(conn, router)
	}
}


func handleRequest(conn net.Conn, router *Router) {
	defer conn.Close()

	var fullData []byte
	buffer := make([]byte, 2048)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			if err.Error() == "EOF" || strings.Contains(err.Error(), "use of closed network connection") {
				break
			}
			fmt.Println("Error reading data:", err)
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			return
		}
		fullData = append(fullData, buffer[:n]...)

		if bytes.Contains(fullData, []byte("\r\n\r\n")) {
			break
		}
		if len(fullData) > 8192 { // 8KB limit for headers
			fmt.Println("Request headers too large")
			conn.Write([]byte("HTTP/1.1 413 Payload Too Large\r\nContent-Length: 22\r\n\r\n413 Payload Too Large"))
			return
		}
	}

	if len(fullData) == 0 {
		fmt.Println("Received empty request data.")
		return
	}

	headerEndIndex := bytes.Index(fullData, []byte("\r\n\r\n"))
	if headerEndIndex == -1 {
		fmt.Println("Invalid HTTP request: Missing \\r\\n\\r\\n")
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 19\r\n\r\n400 Bad Request"))
		return
	}

	requestHeaderBytes := fullData[:headerEndIndex]
	requestBodyBytesSoFar := fullData[headerEndIndex+4:]

	headerLines := bytes.Split(requestHeaderBytes, []byte("\r\n"))
	if len(headerLines) < 1 {
		fmt.Println("Invalid request: No request line found.")
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 19\r\n\r\n400 Bad Request"))
		return
	}

	requestLineParts := bytes.Split(headerLines[0], []byte(" "))
	if len(requestLineParts) < 3 {
		fmt.Println("Invalid request line")
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 19\r\n\r\n400 Bad Request"))
		return
	}

	method := string(bytes.TrimSpace(requestLineParts[0]))
	rawPath := string(bytes.TrimSpace(requestLineParts[1]))
	protocol := string(bytes.TrimSpace(requestLineParts[2]))

	parsedURL, err := url.Parse(rawPath)
	if err != nil {
		fmt.Println("Invalid request path (cannot parse URL):", rawPath)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\nContent-Length: 19\r\n\r\n400 Bad Request"))
		return
	}
	if len(parsedURL.Path) > 2048 { // path length limit
		fmt.Println("Route path too long:", parsedURL.Path)
		conn.Write([]byte("HTTP/1.1 414 URI Too Long\r\nContent-Length: 18\r\n\r\n414 URI Too Long"))
		return
	}

	req := &http.Request{
		Method: method,
		URL:    parsedURL,
		Proto:      protocol,
		Header:     make(http.Header),
		RequestURI: rawPath,
		Body:       http.NoBody,
	}

	for _, line := range headerLines[1:] {
		if len(line) == 0 {
			continue
		}
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) == 2 {
			headerName := strings.TrimSpace(string(parts[0]))
			headerValue := strings.TrimSpace(string(parts[1]))
			req.Header.Add(headerName, headerValue)
			if strings.ToLower(headerName) == "host" {
				req.Host = headerValue
			}
		}
	}

	contentLengthStr := req.Header.Get("Content-Length")
	if contentLengthStr != "" {
		contentLength, err := strconv.ParseInt(contentLengthStr, 10, 64)
		if err == nil && contentLength > 0 {
			
			body := make([]byte, contentLength)
			copiedBytes := copy(body, requestBodyBytesSoFar)

			if int64(copiedBytes) < contentLength {
				bytesToRead := int(contentLength) - copiedBytes

				remainingBodyBuffer := make([]byte, bytesToRead)
				n, readErr := io.ReadFull(conn, remainingBodyBuffer)
				if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
					fmt.Println("Error reading remaining request body:", readErr, "read", n, "expected", bytesToRead)
				}
				copy(body[copiedBytes:], remainingBodyBuffer[:n])
			}
			req.Body = io.NopCloser(bytes.NewReader(body))
		} else if err != nil {
			fmt.Println("Invalid Content-Length:", contentLengthStr)
		}
	}

	resWriter := newConnResponseWriter(conn)
	router.ServeHTTP(resWriter, req)
}
