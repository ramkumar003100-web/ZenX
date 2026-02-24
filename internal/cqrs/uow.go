package cqrs

import "context"

type UnitOfWork interface {
	Do(context.Context, func(context.Context) error) error
}

type FuncUoW struct {
	Exec func(context.Context, func(context.Context) error) error
}

func (u FuncUoW) Do(ctx context.Context, fn func(context.Context) error) error {
	return u.Exec(ctx, fn)
}
