package main

import (
	"fmt"
	"net"
	"os"

	"github.com/radhe5hyam/GoFiber/http"
)

func main() {
	http.RegisterRoute("GET", "/hello", func(params map[string]string) {
		fmt.Println("GET /hello handler called")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := ":" + port

	
	
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error setting up listener:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Server started on port " + port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go http.HandleConnection(conn)
	}
}
