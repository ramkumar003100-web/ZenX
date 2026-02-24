package gateway

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func LoggingInterceptor(l *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		res, err := handler(ctx, req)
		l.Info("grpc request", "method", info.FullMethod, "duration", time.Since(start), "error", err)
		return res, err
	}
}
