package http

import "strings"

type Node struct {
	segment string
	children   map[string]*Node
	paramChild *Node
	paramName string
	isLeaf  bool
	handlers map[string]HandlerFunc
}

type HandlerFunc func(ctx *Context)

func NewNode(segment string) *Node {
	return &Node{
		segment: segment,
		children:  make(map[string]*Node),
		handlers: make(map[string]HandlerFunc),
	}
}

func (n *Node) AddSegment(method string, path string, handler HandlerFunc) {
	segments := splitPath(path)
	current := n

	for _, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			if current.paramChild == nil {
				child := NewNode(segment)
				child.paramName = segment[1:]
				current.paramChild = child
			}
			current = current.paramChild
		} else {
			if _, exist := current.children[segment]; !exist {
				current.children[segment] = NewNode(segment)
			}
			current = current.children[segment]
		}
	}

	current.isLeaf = true
	current.handlers[method] = handler
}

func (n *Node) FindSegment(method string, path string) (HandlerFunc, map[string]string, bool, bool) {
	segments := splitPath(path)
	current := n
	params := make(map[string]string)


	for _, segment := range segments {
		if next, exit := current.children[segment]; exit {
			current = next
		} else if current.paramChild != nil {
			current = current.paramChild
			params[current.paramName] = segment
		} else {
			return nil, nil, false, false
		}
	}

	handler, ok := current.handlers[method]
	if !ok {
		if len(current.handlers) > 0 {
			return nil, nil, true, false
		}
		return nil, nil, false, false
	}
	return handler, params, true, true
}

func splitPath(path string) []string {
	path = strings.Trim(path, " /")
	if path == "" {
		return []string{}
	}
	return strings.Split(path, "/")
}

type Router struct {
	root *Node
}

func NewRouter() *Router {
	return &Router{
		root: NewNode(""),
	}
}

func (r *Router) Register(method string, path string, handler HandlerFunc) {
	r.root.AddSegment(method, path, handler)
}