package runtime

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RecoveryPolicy interface {
	ShouldRestart(module string, err error) bool
	Backoff(module string) time.Duration
}

type ExponentialRecovery struct {
	mu       sync.Mutex
	attempts map[string]int
	max      int
}

func NewExponentialRecovery(max int) *ExponentialRecovery {
	return &ExponentialRecovery{attempts: map[string]int{}, max: max}
}

func (r *ExponentialRecovery) ShouldRestart(module string, _ error) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.attempts[module]++
	return r.attempts[module] <= r.max
}

func (r *ExponentialRecovery) Backoff(module string) time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := r.attempts[module]
	if n < 1 {
		n = 1
	}
	if n > 8 {
		n = 8
	}
	return time.Duration(1<<uint(n-1)) * time.Second
}

func RunWithRecovery(ctx context.Context, module string, fn func(context.Context) error) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("module %s panicked: %v", module, rec)
		}
	}()
	return fn(ctx)
}
