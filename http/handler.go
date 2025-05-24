package http

import (
	"bufio"
	"fmt"
	"net"
)

func Listen(port string, router *Router) {
	address := ":" + port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error setting up listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on port " + port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go router.HandleConnection(conn)
	}
}

func (r *Router)HandleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := ParseHTTPRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}
	writer := NewResponseWriter(conn)

	handler, params, found, allowed := r.root.FindSegment(req.Method, req.Path)
	ctx := &Context{
		Request:  req,
		Response: writer,
		Params:   params,
	}

	switch {
	case !found:
		writer.WriteHeader(404)
		writer.Write([]byte("Not Found"))
	case !allowed:
		writer.WriteHeader(405)
		writer.Write([]byte("Method Not Allowed"))
	case handler != nil:
		handler(ctx)
	default:
		writer.WriteHeader(500)
		writer.Write([]byte("Internal Server Error: No handler found"))
	}
	writer.Flush()
}