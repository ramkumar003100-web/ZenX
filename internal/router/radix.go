package router

import (
	"net/http"
	"strings"
	"sync"
)

type Router struct {
	mu          sync.RWMutex
	trees       map[string]*node
	middlewares []Middleware
}

type node struct {
	segment   string
	handler   HandlerFunc
	children  map[string]*node
	param     *node
	paramName string
}

func New() *Router {
	return &Router{trees: map[string]*node{}}
}

func (r *Router) Use(mw ...Middleware) {
	r.middlewares = append(r.middlewares, mw...)
}

func (r *Router) Handle(method, path string, h HandlerFunc, m ...Middleware) {
	r.mu.Lock()
	defer r.mu.Unlock()

	method = strings.ToUpper(method)
	if r.trees[method] == nil {
		r.trees[method] = &node{children: map[string]*node{}}
	}

	full := Chain(h, append(r.middlewares, m...)...)
	cur := r.trees[method]
	for _, segment := range splitPath(path) {
		if strings.HasPrefix(segment, ":") {
			if cur.param == nil {
				cur.param = &node{children: map[string]*node{}, paramName: strings.TrimPrefix(segment, ":")}
			}
			cur = cur.param
			continue
		}
		if cur.children[segment] == nil {
			cur.children[segment] = &node{segment: segment, children: map[string]*node{}}
		}
		cur = cur.children[segment]
	}
	cur.handler = full
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mu.RLock()
	root := r.trees[req.Method]
	r.mu.RUnlock()

	if root == nil {
		if r.pathExists(req.URL.Path) {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		http.NotFound(w, req)
		return
	}

	ctx := NewContext(w, req)
	cur := root
	for _, seg := range splitPath(req.URL.Path) {
		if next := cur.children[seg]; next != nil {
			cur = next
			continue
		}
		if cur.param != nil {
			ctx.SetParam(cur.param.paramName, seg)
			cur = cur.param
			continue
		}
		http.NotFound(w, req)
		return
	}

	if cur.handler == nil {
		http.NotFound(w, req)
		return
	}
	cur.handler(ctx)
}

func (r *Router) pathExists(path string) bool {
	for _, root := range r.trees {
		cur := root
		matched := true
		for _, seg := range splitPath(path) {
			if next := cur.children[seg]; next != nil {
				cur = next
				continue
			}
			if cur.param != nil {
				cur = cur.param
				continue
			}
			matched = false
			break
		}
		if matched && cur.handler != nil {
			return true
		}
	}
	return false
}

func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "/")
}
