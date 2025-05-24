package http

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)


type ResponseWriter struct {
	conn net.Conn
	writer *bufio.Writer
	statusCode int
	headers map[string]string
	body []byte
	flushed bool
}


func NewResponseWriter(conn net.Conn) *ResponseWriter{
	return &ResponseWriter{
		conn : conn,
		writer: bufio.NewWriter(conn),
		headers: make(map[string]string),

	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) Header(key, value string) {
	rw.headers[key] = value
}

func (rw *ResponseWriter) Write(data []byte) {
	rw.body = append(rw.body, data...)
}

func (rw *ResponseWriter) Flush() error {
	if rw.flushed {
		return nil
	}
	rw.flushed = true

	if rw.statusCode == 0 {
		rw.statusCode = 200
	}

	
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", rw.statusCode, statusText(rw.statusCode))
	rw.writer.WriteString(statusLine)

	if _, ok := rw.headers["Content-Length"]; !ok {
		rw.headers["Content-Length"] = strconv.Itoa(len(rw.body))
	}
	if _, ok := rw.headers["Content-Type"]; !ok {
		rw.headers["Content-Type"] = "text/plain"
	}

	for key, value := range rw.headers {
		rw.writer.WriteString(fmt.Sprintf(("%s: %s\r\n"), key,value))
	}

	rw.writer.WriteString("\r\n")
	rw.writer.Write(rw.body)

	return rw.writer.Flush()
}

func statusText(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}