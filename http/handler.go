package http

import (
	"net"

	"github.com/radhe5hyam/GoFiber/http/router"
)

var routeRoot = router.NewNode("/")

func RegisterRoute(method, path string, handler router.HandlerFunc) {
	routeRoot.AddSegment(method, path, handler)
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := ParseHTTPRequest(conn)
	writer := ResponseWriter(conn)

	if err != nil {
		writer.WriteHeader(400)
		writer.body.WriteString("Bad Request: " + err.Error())
		writer.Flush()
		return
	}

	handler, params, found, allowed := routeRoot.FindSegment(req.Method, req.Path)
	switch {
	case !found:
		writer.WriteHeader(404)
		writer.body.WriteString("Not Found")
	case !allowed:
		writer.WriteHeader(405)
		writer.body.WriteString("Method Not Allowed")
	case handler != nil:
		handler(params)
	default:
		writer.WriteHeader(500)
		writer.body.WriteString("Internal Server Error: No handler found")
	}
	writer.Flush()
}