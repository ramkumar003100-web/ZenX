package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"zenx/internal/router"
)

type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
}

func RateLimit(rps float64, burst int) router.Middleware {
	var mu sync.Mutex
	clients := map[string]*tokenBucket{}

	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
			if ip == "" {
				ip = c.Request.RemoteAddr
			}
			now := time.Now()

			mu.Lock()
			bucket, ok := clients[ip]
			if !ok {
				bucket = &tokenBucket{tokens: float64(burst), lastRefill: now}
				clients[ip] = bucket
			}
			elapsed := now.Sub(bucket.lastRefill).Seconds()
			bucket.tokens += elapsed * rps
			if bucket.tokens > float64(burst) {
				bucket.tokens = float64(burst)
			}
			bucket.lastRefill = now

			if bucket.tokens < 1 {
				mu.Unlock()
				http.Error(c.Writer, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			bucket.tokens--
			mu.Unlock()

			next(c)
		}
	}
}
