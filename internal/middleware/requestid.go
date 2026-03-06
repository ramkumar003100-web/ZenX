package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"zenx/internal/router"
)

type contextKey string

const RequestIDContextKey contextKey = "zenx.request_id"

func RequestID(headerName string) router.Middleware {
	if headerName == "" {
		headerName = "X-Request-ID"
	}

	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			requestID := c.Request.Header.Get(headerName)
			if requestID == "" {
				requestID = newRequestID()
			}

			ctx := context.WithValue(c.Request.Context(), RequestIDContextKey, requestID)
			c.WithContext(ctx)
			c.Writer.Header().Set(headerName, requestID)

			next(c)
		}
	}
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(RequestIDContextKey).(string); ok {
		return v
	}
	return ""
}

func newRequestID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "generr"
	}
	return hex.EncodeToString(b[:])
}

func IsSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
