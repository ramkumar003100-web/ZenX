package gateway

import (
	"context"
	"sync"
)

type CollapseGroup struct {
	mu       sync.Mutex
	inFlight map[string]*call
}
type call struct {
	done chan struct{}
	val  any
	err  error
}

func NewCollapseGroup() *CollapseGroup { return &CollapseGroup{inFlight: map[string]*call{}} }
func (g *CollapseGroup) Do(ctx context.Context, key string, fn func(context.Context) (any, error)) (any, error) {
	g.mu.Lock()
	if c, ok := g.inFlight[key]; ok {
		g.mu.Unlock()
		<-c.done
		return c.val, c.err
	}
	c := &call{done: make(chan struct{})}
	g.inFlight[key] = c
	g.mu.Unlock()
	c.val, c.err = fn(ctx)
	close(c.done)
	g.mu.Lock()
	delete(g.inFlight, key)
	g.mu.Unlock()
	return c.val, c.err
}
