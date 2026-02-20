package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"zenx/internal/router"
	"zenx/pkg/logger"
)

type timeoutWriter struct {
	h      http.Header
	body   []byte
	status int
	w      http.ResponseWriter
	mu     sync.Mutex
}

func newTimeoutWriter(w http.ResponseWriter) *timeoutWriter {
	return &timeoutWriter{h: make(http.Header), w: w, status: http.StatusOK}
}

func (tw *timeoutWriter) Header() http.Header { return tw.h }

func (tw *timeoutWriter) WriteHeader(statusCode int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.status = statusCode
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	tw.body = append(tw.body, b...)
	return len(b), nil
}

func (tw *timeoutWriter) flush() {
	for k, values := range tw.h {
		for _, v := range values {
			tw.w.Header().Add(k, v)
		}
	}
	tw.w.WriteHeader(tw.status)
	if len(tw.body) > 0 {
		_, _ = tw.w.Write(tw.body)
	}
}

func Timeout(d time.Duration) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			ctx, cancel := context.WithTimeout(c.Request.Context(), d)
			defer cancel()

			tw := newTimeoutWriter(c.Writer)
			cc := *c
			cc.Writer = tw
			cc.WithContext(ctx)

			done := make(chan struct{})
			go func() {
				defer close(done)
				next(&cc)
			}()

			select {
			case <-done:
				tw.flush()
			case <-ctx.Done():
				http.Error(c.Writer, "request timeout", http.StatusGatewayTimeout)
			}
		}
	}
}

func RequestID() router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			requestID := c.Request.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = time.Now().UTC().Format("20060102150405.000000000")
			}
			c.Writer.Header().Set("X-Request-ID", requestID)
			ctx := logger.WithRequestID(c.Request.Context(), requestID)
			c.WithContext(ctx)
			next(c)
		}
	}
}