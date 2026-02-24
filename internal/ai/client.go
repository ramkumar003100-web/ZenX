package ai

import "context"

type Client interface {
	Embed(context.Context, string) ([]float32, error)
	Complete(context.Context, string) (string, error)
}
