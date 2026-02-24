package performance

import (
	"sync/atomic"
	"time"
)

type AdaptiveLimiter struct {
	limit    atomic.Int64
	min, max int64
}

func NewAdaptiveLimiter(min, max int64) *AdaptiveLimiter {
	a := &AdaptiveLimiter{min: min, max: max}
	a.limit.Store(min)
	return a
}
func (a *AdaptiveLimiter) Adjust(latency time.Duration) {
	cur := a.limit.Load()
	if latency > 300*time.Millisecond && cur > a.min {
		a.limit.Store(cur - 1)
	} else if latency < 100*time.Millisecond && cur < a.max {
		a.limit.Store(cur + 1)
	}
}
func (a *AdaptiveLimiter) Limit() int64 { return a.limit.Load() }
