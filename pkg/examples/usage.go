package examples

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"zenx/internal/featureflags"
	"zenx/internal/jobs"
	"zenx/internal/notifications"
	"zenx/internal/storage"
	"zenx/internal/tracing"
	"zenx/internal/websocket"
)

func WebSocketExample(h *websocket.Hub) {
	_ = h.SendToClient("u-1", []byte(`{"event":"hello"}`))
	h.Broadcast(context.Background(), []byte("system maintenance"))
}

func DistributedJobsExample(q *jobs.DistributedQueue) {
	_ = q.Enqueue(context.Background(), "send-email", map[string]any{"to": "a@b.com"}, time.Second)
	handlers := map[string]func(context.Context, json.RawMessage) error{
		"send-email": func(ctx context.Context, msg json.RawMessage) error { return nil },
	}
	q.Consume(context.Background(), handlers, time.Second)
}

func FeatureFlagExample(ff *featureflags.Manager) bool {
	ff.Set("new-dashboard", featureflags.FlagRule{Roles: map[string]bool{"admin": true}})
	return ff.Enabled("new-dashboard", "u-1", "admin")
}

func EmailExample(sender *notifications.EmailSender) error {
	return sender.SendTemplate(context.Background(), "dev@example.com", "ZenX", "Hello {{.Name}}", map[string]string{"Name": "Team"}, 2)
}

func FileUploadExample(handler http.HandlerFunc, s *storage.Service) http.HandlerFunc {
	_ = s
	return handler
}

func TracingExample(t *tracing.Tracer) {
	ctx, span := t.Start(context.Background(), "example-span")
	defer span.End()
	_ = ctx
}

func IntegratedFlowExample() {
	// Example flow: validate request -> authenticate -> cache -> enqueue job.
	// Implementors should wire these with router middlewares in real handlers.
}
