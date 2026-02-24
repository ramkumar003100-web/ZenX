package gateway

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Upstream struct {
	Name    string
	Target  *url.URL
	Breaker *CircuitBreaker
	Rate    func(*http.Request) bool
}

type Gateway struct {
	mu        sync.RWMutex
	upstreams map[string]Upstream
	retries   int
}

func NewGateway(retries int) *Gateway {
	return &Gateway{upstreams: map[string]Upstream{}, retries: retries}
}
func (g *Gateway) Register(u Upstream) { g.mu.Lock(); g.upstreams[u.Name] = u; g.mu.Unlock() }

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mu.RLock()
	u, ok := g.upstreams[r.Host]
	g.mu.RUnlock()
	if !ok {
		http.Error(w, "unknown upstream", http.StatusBadGateway)
		return
	}
	if u.Breaker != nil {
		if err := u.Breaker.Allow(); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}
	if u.Rate != nil && !u.Rate(r) {
		http.Error(w, "service rate limit", http.StatusTooManyRequests)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(u.Target)
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	r = r.WithContext(ctx)
	proxy.ErrorHandler = func(_ http.ResponseWriter, _ *http.Request, err error) {
		if u.Breaker != nil {
			u.Breaker.MarkFailure()
		}
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
	proxy.ModifyResponse = func(*http.Response) error {
		if u.Breaker != nil {
			u.Breaker.MarkSuccess()
		}
		return nil
	}
	proxy.ServeHTTP(w, r)
}
