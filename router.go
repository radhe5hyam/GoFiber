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

type Router struct{
	Root *Node
}

func NewRouter() *Router {
    return &Router{
        Root: NewNode(),
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

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	cleanedPath := path.Clean(req.URL.Path)
    segments := splitPath(cleanedPath)
    current := r.Root

    for _, segment := range segments {
        if nextNode, ok := current.Children[segment]; ok {
            current = nextNode
        } else {
			// Path segment not found
            http.NotFound(res, req)
            return
        }
    }

	// Path found, now check if a handler exists for the request method
    if handler, ok := current.Handlers[req.Method]; ok {
        handler(res, req) 
    } else if len(current.Handlers) > 0 {
		allowedMethods := make([]string, 0, len(current.Handlers))
		for m := range current.Handlers {
			allowedMethods = append(allowedMethods, m)
		}
		res.Header().Set("Allow", strings.Join(allowedMethods, ", "))
		http.Error(res, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    } else {
		// Path is valid, but no handlers are registered for any method on this node
        http.NotFound(res, req)
    }
}

func splitPath(path string) []string {
    return strings.Split(strings.Trim(path, "/"), "/")
}