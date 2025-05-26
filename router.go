package main

import (
	"net/http"
	"path"
	"strings"
)

type Node struct {
	Children map[string]*Node
	Handlers map[string]http.HandlerFunc
}

func NewNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
		Handlers: make(map[string]http.HandlerFunc),
	}
}

// MiddlewareFunc defines the signature for middleware functions.
type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

type Router struct{
	Root       *Node
	Middleware []MiddlewareFunc // Slice to store global middleware
}

func NewRouter() *Router {
    return &Router{
        Root:       NewNode(),
        Middleware: make([]MiddlewareFunc, 0),
    }
}

func (r *Router) AddRoute(method, routePath string, handler http.HandlerFunc) {
	cleanedPath := path.Clean(routePath)
	current := r.Root
	parts := splitPath(cleanedPath)

	for _, part := range parts {
		if _, exists := current.Children[part]; !exists {
			current.Children[part] = NewNode()
		}
		current = current.Children[part]
	}

	current.Handlers[method] = handler
}

// Use adds a new global middleware to the router.
func (r *Router) Use(middleware MiddlewareFunc) {
	r.Middleware = append(r.Middleware, middleware)
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	cleanedPath := path.Clean(req.URL.Path)
    segments := splitPath(cleanedPath)
    current := r.Root
	var targetHandler http.HandlerFunc

    for _, segment := range segments {
        if nextNode, ok := current.Children[segment]; ok {
            current = nextNode
        } else {
			targetHandler = http.NotFound
			// Apply middleware even for NotFound
			for i := len(r.Middleware) - 1; i >= 0; i-- {
				targetHandler = r.Middleware[i](targetHandler)
			}
			targetHandler(res, req)
            return
        }
    }

    if handler, ok := current.Handlers[req.Method]; ok {
		targetHandler = handler
    } else if len(current.Handlers) > 0 {
		var allowedMethods []string
		for m := range current.Handlers {
			allowedMethods = append(allowedMethods, m)
		}
		res.Header().Set("Allow", strings.Join(allowedMethods, ", "))
		targetHandler = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
    } else {
		targetHandler = http.NotFound
    }

	// Apply all global middleware
	for i := len(r.Middleware) - 1; i >= 0; i-- {
		targetHandler = r.Middleware[i](targetHandler)
	}
	targetHandler(res, req)
}

func splitPath(path string) []string {
    return strings.Split(strings.Trim(path, "/"), "/")
}