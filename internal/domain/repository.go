package domain

import "context"

type Repository[T AggregateRoot] interface {
	Get(context.Context, string) (T, error)
	Save(context.Context, T) error
}
