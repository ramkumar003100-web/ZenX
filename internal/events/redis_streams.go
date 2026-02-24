package events

import "context"

type RedisStreams interface {
	XAdd(context.Context, string, map[string]any) error
	XRead(context.Context, string, string) (map[string]any, error)
}

type RedisStreamsAdapter struct{ Client RedisStreams }

func (a RedisStreamsAdapter) Publish(ctx context.Context, stream string, e Event) error {
	return a.Client.XAdd(ctx, stream, map[string]any{"name": e.Name, "payload": e.Payload})
}
