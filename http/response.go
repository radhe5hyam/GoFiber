package http

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)


type HTTPResponse struct {
	conn net.Conn
	statusCode int
	headers map[string]string
	body strings.Builder
	writer *bufio.Writer
	wrote bool
}


func ResponseWriter(conn net.Conn) *HTTPResponse{
	return &HTTPResponse{
		conn : conn,
		writer: bufio.NewWriter(conn),
		statusCode: 200,
		headers: make(map[string]string),

	}
}

func (w *HTTPResponse) WriteHeader(statusCode int) {
	if w.wrote {
		return
	}
	w.statusCode = statusCode
	w.wrote = true
}

func (w *HTTPResponse) Header(key, value string) {
	w.headers[key] = value
}

func (w *HTTPResponse) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func (w *HTTPResponse) Flush() error {
	bodyStr := w.body.String()

	w.headers["Content-length"] = strconv.Itoa(len(bodyStr))
	if _, ok := w.headers["Content-Type"]; !ok {
		w.headers["Content-Type"] = "text/plain"
	}

	statusText := statusText(w.statusCode)
	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", w.statusCode, statusText)

	if _, err := w.writer.WriteString((statusLine)); err != nil {
		return err
	}

	for key, value := range w.headers {
		headerLine := fmt.Sprintf(("%s: %s\r\n"), key,value)
		if _, err := w.writer.WriteString(headerLine); err != nil {
			return err
		}
	}

	if _, err := w.writer.WriteString("\r\n"); err != nil {
		return err
	}

	if _, err := w.writer.WriteString(bodyStr); err != nil {
		return err
	}

	return w.writer.Flush()
}

func statusText(statusCode int) string {
	switch statusCode {
	case 200:
		return "OK"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}