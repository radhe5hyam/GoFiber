package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

func ParseHTTPRequest(conn net.Conn) (*HTTPRequest, error) {
	reader := bufio.NewReader(conn)

	// requetst line
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	
	requestLine  = strings.TrimSpace(requestLine)
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid request line: %s", requestLine)
	}

	method := parts[0]
	path := parts[1]


	// headers
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("error reading headers: %v", err)
		}
		if strings.TrimSpace(line) == "" {
			break
		}

		line = strings.TrimSpace(line)
		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			headers[headerParts[0]] = strings.TrimSpace(headerParts[1])
		}
	}
	
	// body
	body := ""
	bodyLength, ok := headers["Content-Length"];
	length, err := strconv.Atoi(bodyLength)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Length header: %v", err)
	}
	if ok {
		buf := make([]byte, length)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, fmt.Errorf("error reading body: %v", err)
		}
		body = string(buf)
	}




	return &HTTPRequest{
		Method:  method,
		Path:    path,
		Headers: headers,
		Body:    body,
	}, nil
}