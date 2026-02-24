package ai

import "context"

type VectorDB interface {
	Upsert(context.Context, string, []float32, map[string]string) error
	Query(context.Context, []float32, int) ([]string, error)
}
