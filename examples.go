package main

import (
	"fmt"
	"net/http"
)


func SetupExampleRoutes(router *Router) {
	router.AddRoute(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "RootPath")
		fmt.Fprint(w, "Hello from the root path, served by the router!")
	})

	router.AddRoute(http.MethodGet, "/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world from /hello path, served by the router!")
	})

	router.AddRoute(http.MethodGet, "/user/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"name": "Router", "type": "Trie"}`)
	})

	router.AddRoute(http.MethodPost, "/submit", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "Data submission acknowledged by router!")
	})
}