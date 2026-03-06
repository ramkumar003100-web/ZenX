package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"zenx/internal/router"
)

func Recovery(log *slog.Logger) router.Middleware {
	if log == nil {
		log = slog.Default()
	}

	return func(next router.HandlerFunc) router.HandlerFunc {
		return func(c *router.Context) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered",
						slog.Any("panic", rec),
						slog.String("path", c.Request.URL.Path),
						slog.String("method", c.Request.Method),
						slog.String("request_id", GetRequestID(c.Request.Context())),
						slog.String("stack", string(debug.Stack())),
					)
					http.Error(c.Writer, fmt.Sprintf("internal server error"), http.StatusInternalServerError)
				}
			}()
			next(c)
		}
	}
}
