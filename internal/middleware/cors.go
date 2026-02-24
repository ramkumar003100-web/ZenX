package middleware

import (
	"net/http"
	"strings"

	"zenx/internal/router"
)

func CORS(allowedOrigins []string) router.Middleware {
	allowed := map[string]struct{}{}
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}

	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			origin := c.Request.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowed["*"]; ok || containsOrigin(allowed, origin) {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					c.Writer.Header().Set("Vary", "Origin")
				}
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			}
			if c.Request.Method == http.MethodOptions {
				c.Writer.WriteHeader(http.StatusNoContent)
				return
			}
			next(c)
		}
	}
}

func containsOrigin(allowed map[string]struct{}, origin string) bool {
	for o := range allowed {
		if strings.EqualFold(o, origin) {
			return true
		}
	}
	return false
}
