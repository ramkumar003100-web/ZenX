package events

import (
	"context"
	"sync"
)

type Event struct {
	Name    string
	Payload any
	Headers map[string]string
}
type Middleware func(context.Context, Event, Handler) error
type Handler func(context.Context, Event) error

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
	mw       []Middleware
	dlq      Handler
}

func NewBus() *Bus              { return &Bus{handlers: map[string][]Handler{}} }
func (b *Bus) Use(m Middleware) { b.mw = append(b.mw, m) }
func (b *Bus) On(name string, h Handler) {
	b.mu.Lock()
	b.handlers[name] = append(b.handlers[name], h)
	b.mu.Unlock()
}
func (b *Bus) WithDLQ(h Handler) { b.dlq = h }
func (b *Bus) Publish(ctx context.Context, e Event) {
	b.mu.RLock()
	hs := append([]Handler{}, b.handlers[e.Name]...)
	b.mu.RUnlock()
	for _, h := range hs {
		handler := h
		for i := len(b.mw) - 1; i >= 0; i-- {
			m := b.mw[i]
			next := handler
			handler = func(c context.Context, ev Event) error { return m(c, ev, next) }
		}
		if err := handler(ctx, e); err != nil && b.dlq != nil {
			_ = b.dlq(ctx, e)
		}
	}
}
