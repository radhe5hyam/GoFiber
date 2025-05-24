package http

type Context struct {
	Request  *Request
	Response *ResponseWriter
	Params   map[string]string
}