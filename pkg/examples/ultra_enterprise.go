package examples

import (
	"context"
	"net/http"
	"time"

	"zenx/internal/billing"
	"zenx/internal/events"
	"zenx/internal/gateway"
	"zenx/internal/multitenancy"
	"zenx/internal/oauth"
	"zenx/internal/plugins"
)

func KafkaEventExample(adapter events.KafkaAdapter) error {
	return adapter.Publish(context.Background(), "user.created", events.Event{Name: "user.created", Payload: map[string]any{"id": "u-1"}})
}

func GRPCCallExample(ctx context.Context, target string) error {
	conn, err := gateway.NewGRPCClientConn(ctx, target)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func MultiTenantQueryExample(ctx context.Context) (string, error) {
	ctx = multitenancy.WithTenant(ctx, "tenant-a")
	return multitenancy.ScopedQuery(ctx, "SELECT * FROM invoices")
}

func OAuth2LoginExample(server *oauth.OAuth2Server) http.HandlerFunc { return server.TokenHandler() }

func PluginLoadingExample(ctx context.Context, reg *plugins.Registry) error {
	d := plugins.Descriptor{Name: "analytics", Version: "1.0.0", Entry: "./plugins/analytics.wasm"}
	if err := reg.Register(d); err != nil {
		return err
	}
	return plugins.DownloadPlugin(ctx, "https://plugins.example.com/analytics.wasm", "/tmp/analytics.wasm")
}

func StripeBillingEventExample() http.HandlerFunc {
	return billing.StripeWebhook(func(_ context.Context, e billing.StripeEvent) error {
		_ = e
		return nil
	})
}

func RequestCollapsingExample(g *gateway.CollapseGroup) (any, error) {
	return g.Do(context.Background(), "same-request-key", func(context.Context) (any, error) {
		time.Sleep(10 * time.Millisecond)
		return map[string]string{"status": "ok"}, nil
	})
}
