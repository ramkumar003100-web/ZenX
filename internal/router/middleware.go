package router

func Chain(h HandlerFunc, m ...Middleware) HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
