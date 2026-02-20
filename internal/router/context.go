package router

import (
	"context"
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	params  map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{Writer: w, Request: r, params: map[string]string{}}
}

func (c *Context) Param(name string) string { return c.params[name] }
func (c *Context) SetParam(k, v string)     { c.params[k] = v }
func (c *Context) WithContext(ctx context.Context) {
	c.Request = c.Request.WithContext(ctx)
}
