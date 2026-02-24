package kubernetes

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(ctx context.Context, server *http.Server, timeout time.Duration, onShutdown func(context.Context)) error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
	case <-stop:
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if onShutdown != nil {
		onShutdown(shutdownCtx)
	}
	return server.Shutdown(shutdownCtx)
}
