package http

import (
	"bufio"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Query   map[string]string
	Body    []byte
}

func ParseHTTPRequest(r io.Reader) (*Request, error) {
	scanner := bufio.NewScanner(r)
	req := &Request{
		Headers: make(map[string]string),
		Query:   make(map[string]string),
	}

	if !scanner.Scan() {
		return nil, io.EOF
	}

	// requetst line
	parts := strings.Split(scanner.Text(), " ")
	if len(parts) != 3 {
		return nil, io.ErrUnexpectedEOF
	}

	req.Method = parts[0]
	rawURL := parts[1]

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	req.Path = parsedURL.Path

	// parse query parameters
	for key, values := range parsedURL.Query() {
		if len(values) > 0 {
			req.Query[key] = values[0]
		}
	}

	// headers
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		headerParts := strings.SplitN(line, ":", 2)
		if len(headerParts) == 2 {
			req.Headers[strings.ToLower(headerParts[0])] = strings.TrimSpace(headerParts[1])
		}
	}
	
	// body
	contentLength := 0
	if cl, ok := req.Headers["content-length"]; ok {
		contentLength, _ = strconv.Atoi(cl)
	}

	if contentLength <= 0 {
		body := make([]byte, contentLength)
		io.ReadFull(r, body)
		req.Body = body
	}

	return req, nil
}