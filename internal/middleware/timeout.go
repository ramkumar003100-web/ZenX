package middleware

import (
	"context"
	"net/http"
	"time"

	"zenx/internal/router"
	"zenx/pkg/logger"
)

func Timeout(d time.Duration) router.Middleware {
	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			ctx, cancel := context.WithTimeout(c.Request.Context(), d)
			defer cancel()
			c.WithContext(ctx)
			next(c)
			if ctx.Err() == context.DeadlineExceeded {
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
