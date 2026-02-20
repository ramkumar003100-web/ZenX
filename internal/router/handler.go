package router

import "net/http"

type HandlerFunc func(*Context)
type Middleware func(HandlerFunc) HandlerFunc

func WrapHTTP(h http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		h(c.Writer, c.Request)
	}
}
