package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func StartSpan(ctx context.Context, tracerName, name string) (context.Context, trace.Span) {
	return otel.Tracer(tracerName).Start(ctx, name)
}
