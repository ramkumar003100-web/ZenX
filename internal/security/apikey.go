package security

import (
	"crypto/subtle"
	"net/http"

	"zenx/internal/router"
)

func APIKeyAuth(expected string) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			key := c.Request.Header.Get("X-API-Key")
			if subtle.ConstantTimeCompare([]byte(key), []byte(expected)) != 1 {
				http.Error(c.Writer, "invalid api key", http.StatusUnauthorized)
				return
			}
			next(c)
		}
	}
}
