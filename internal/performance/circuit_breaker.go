package performance

import (
	"errors"
	"sync"
	"time"
)

type Breaker struct {
	mu        sync.Mutex
	open      bool
	halfOpen  bool
	fail      int
	threshold int
	opened    time.Time
	cool      time.Duration
}

func NewBreaker(threshold int, cool time.Duration) *Breaker {
	return &Breaker{threshold: threshold, cool: cool}
}
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.open {
		if time.Since(b.opened) >= b.cool {
			b.halfOpen = true
			b.open = false
			return nil
		}
		return errors.New("breaker open")
	}
	return nil
}
func (b *Breaker) Success() { b.mu.Lock(); b.fail = 0; b.halfOpen = false; b.mu.Unlock() }
func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.fail++
	if b.fail >= b.threshold || b.halfOpen {
		b.open = true
		b.opened = time.Now()
	}
}
