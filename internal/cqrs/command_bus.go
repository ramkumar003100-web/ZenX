package cqrs

import (
	"context"
	"fmt"
	"sync"
)

type Command interface{ Name() string }
type CommandHandler func(context.Context, Command) error

type CommandBus struct {
	mu       sync.RWMutex
	handlers map[string]CommandHandler
}

func NewCommandBus() *CommandBus { return &CommandBus{handlers: map[string]CommandHandler{}} }
func (b *CommandBus) Register(name string, h CommandHandler) {
	b.mu.Lock()
	b.handlers[name] = h
	b.mu.Unlock()
}
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	b.mu.RLock()
	h, ok := b.handlers[cmd.Name()]
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("no handler for %s", cmd.Name())
	}
	return h(ctx, cmd)
}
