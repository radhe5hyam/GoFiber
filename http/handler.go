package http

import (
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := ParseHTTPRequest(conn)
	if err != nil {
		response := "HTTP/1.1 400 Bad Request\r\nContent-Type: text/plain\r\n\r\n" + err.Error()
		conn.Write(([]byte(response)))
		return
	}
	fmt.Println("Received request:", req.Method, req.Path, req.Body)

	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
	conn.Write([]byte(response))
}