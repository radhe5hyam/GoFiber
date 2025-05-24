package main

import (
	"os"

	"github.com/radhe5hyam/GoFiber/http"
)

func main() {
	router := http.NewRouter()
	RegisterRoutes(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.Listen(port, router)
}

func RegisterRoutes(router *http.Router) {
	router.Register("GET", "/", func(ctx *http.Context) {
		ctx.Response.Header("Content-Type", "text/plain")
		ctx.Response.WriteHeader(200)
		ctx.Response.Write([]byte("It Works!"))
		ctx.Response.Flush()
	})
}

