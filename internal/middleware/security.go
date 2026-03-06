package middleware

import (
	"net/http"
	"sync"
	"time"

	"zenx/internal/router"
)

func SecureHeaders() router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			h := c.Writer.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("X-XSS-Protection", "1; mode=block")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			next(c)
		}
	}
}

func CSRF(headerName string) router.Middleware {
	if headerName == "" {
		headerName = "X-CSRF-Token"
	}
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			if IsSafeMethod(c.Request.Method) {
				next(c)
				return
			}
			if c.Request.Header.Get(headerName) == "" {
				http.Error(c.Writer, "missing csrf token", http.StatusForbidden)
				return
			}
			next(c)
		}
	}
}

func BruteForceGuard(maxAttempts int, period time.Duration) router.Middleware {
	type entry struct {
		n int
		t time.Time
	}
	attempts := map[string]entry{}
	var mu sync.Mutex

	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			key := clientIP(c.Request)
			mu.Lock()
			e := attempts[key]
			if time.Since(e.t) > period {
				e = entry{}
			}
			if e.n >= maxAttempts {
				mu.Unlock()
				http.Error(c.Writer, "too many attempts", http.StatusTooManyRequests)
				return
			}
			e.n++
			e.t = time.Now()
			attempts[key] = e
			mu.Unlock()
			next(c)
		}
	}
}
