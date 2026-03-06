package router

import "strings"

type Group struct {
	router      *Router
	prefix      string
	middlewares []Middleware
}

func (r *Router) Group(prefix string, m ...Middleware) *Group {
	return &Group{router: r, prefix: normalizeGroupPrefix(prefix), middlewares: m}
}

func (g *Group) Group(prefix string, m ...Middleware) *Group {
	combined := make([]Middleware, 0, len(g.middlewares)+len(m))
	combined = append(combined, g.middlewares...)
	combined = append(combined, m...)

	return &Group{
		router:      g.router,
		prefix:      joinPaths(g.prefix, prefix),
		middlewares: combined,
	}
}

func (g *Group) Handle(method, path string, h HandlerFunc, m ...Middleware) {
	all := make([]Middleware, 0, len(g.middlewares)+len(m))
	all = append(all, g.middlewares...)
	all = append(all, m...)
	g.router.Handle(method, joinPaths(g.prefix, path), h, all...)
}

func (g *Group) Version(version string) *Group {
	v := strings.TrimSpace(version)
	if v == "" {
		return g
	}
	if !strings.HasPrefix(v, "v") {
		v = "v" + v
	}
	return g.Group("/" + v)
}

func normalizeGroupPrefix(prefix string) string {
	p := strings.TrimSpace(prefix)
	if p == "" || p == "/" {
		return ""
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.TrimSuffix(p, "/")
}

func joinPaths(prefix, path string) string {
	base := normalizeGroupPrefix(prefix)
	suffix := strings.TrimSpace(path)
	if suffix == "" || suffix == "/" {
		if base == "" {
			return "/"
		}
		return base
	}
	if !strings.HasPrefix(suffix, "/") {
		suffix = "/" + suffix
	}
	if base == "" {
		return suffix
	}
	return base + suffix
}
