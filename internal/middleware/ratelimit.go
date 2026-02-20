package middleware

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"zenx/internal/router"
)

type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
	lastSeen   time.Time
}

func RateLimit(rps float64, burst int) router.Middleware {
	var mu sync.Mutex
	clients := map[string]*tokenBucket{}
	cleanupTicker := time.NewTicker(5 * time.Minute)
	go func() {
		for range cleanupTicker.C {
			now := time.Now()
			mu.Lock()
			for ip, bucket := range clients {
				if now.Sub(bucket.lastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			ip := clientIP(c.Request)
			now := time.Now()

			mu.Lock()
			bucket, ok := clients[ip]
			if !ok {
				bucket = &tokenBucket{tokens: float64(burst), lastRefill: now, lastSeen: now}
				clients[ip] = bucket
			}
			elapsed := now.Sub(bucket.lastRefill).Seconds()
			bucket.tokens += elapsed * rps
			if bucket.tokens > float64(burst) {
				bucket.tokens = float64(burst)
			}
			bucket.lastRefill = now
			bucket.lastSeen = now

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

func clientIP(r *http.Request) string {
	xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	xri := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if xri != "" {
		return xri
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if ip == "" {
		return r.RemoteAddr
	}
	return ip
}
