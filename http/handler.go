package http

import (
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := ParseHTTPRequest(conn)
	if err != nil {
		writer := ResponseWriter(conn)
		writer.WriteHeader(400)
		writer.body.WriteString("Bad Request: " + err.Error())
		writer.Flush()
		return
	}
	fmt.Println("Received request:", req.Method, req.Path, req.Body)

	writer := ResponseWriter(conn)
	writer.WriteHeader(200)
	writer.Header("X-Custom-Header", "GoFiber")
	writer.Write([]byte("It works! - GoFiber"))
	writer.Flush()
}