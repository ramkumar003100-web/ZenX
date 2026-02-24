package cqrs

import (
	"context"
	"fmt"
	"sync"
)

type Query interface{ Name() string }
type QueryHandler func(context.Context, Query) (any, error)

type QueryBus struct {
	mu       sync.RWMutex
	handlers map[string]QueryHandler
}

func NewQueryBus() *QueryBus { return &QueryBus{handlers: map[string]QueryHandler{}} }
func (b *QueryBus) Register(name string, h QueryHandler) {
	b.mu.Lock()
	b.handlers[name] = h
	b.mu.Unlock()
}
func (b *QueryBus) Ask(ctx context.Context, q Query) (any, error) {
	b.mu.RLock()
	h, ok := b.handlers[q.Name()]
	b.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no handler for %s", q.Name())
	}
	return h(ctx, q)
}
