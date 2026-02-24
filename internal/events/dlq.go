package events

import (
	"context"
	"sync"
)

type InMemoryDLQ struct {
	mu    sync.Mutex
	items []Event
}

func (d *InMemoryDLQ) Push(_ context.Context, e Event) error {
	d.mu.Lock()
	d.items = append(d.items, e)
	d.mu.Unlock()
	return nil
}
func (d *InMemoryDLQ) Items() []Event {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]Event, len(d.items))
	copy(out, d.items)
	return out
}
