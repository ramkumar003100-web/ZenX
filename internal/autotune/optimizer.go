package autotune

import (
	"context"
	"sync"
	"time"
)

type Applier interface {
	Apply(context.Context, Recommendation) error
}

type Optimizer struct {
	mu      sync.RWMutex
	current Recommendation
	applier Applier
}

func NewOptimizer(applier Applier) *Optimizer { return &Optimizer{applier: applier} }

func (o *Optimizer) Start(ctx context.Context, collect func(context.Context) RuntimeMetrics, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			metrics := collect(ctx)
			o.mu.Lock()
			o.current = Analyze(metrics, o.current)
			rec := o.current
			o.mu.Unlock()
			_ = o.applier.Apply(ctx, rec)
		}
	}
}

func (o *Optimizer) Current() Recommendation {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.current
}
