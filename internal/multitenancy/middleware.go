package multitenancy

import (
	"net/http"

	"zenx/internal/router"
)

func Resolver(header string) router.Middleware {
	if header == "" {
		header = "X-Tenant-ID"
	}
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			tenant := c.Request.Header.Get(header)
			if tenant == "" {
				http.Error(c.Writer, "missing tenant", http.StatusBadRequest)
				return
			}
			c.WithContext(WithTenant(c.Request.Context(), tenant))
			next(c)
		}
	}
}
