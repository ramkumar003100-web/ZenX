package gateway

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	mu        sync.Mutex
	state     State
	failures  int
	threshold int
	openUntil time.Time
	cooldown  time.Duration
}

func NewCircuitBreaker(threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{state: Closed, threshold: threshold, cooldown: cooldown}
}

func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == Open {
		if time.Now().After(cb.openUntil) {
			cb.state = HalfOpen
			return nil
		}
		return errors.New("circuit open")
	}
	return nil
}
func (cb *CircuitBreaker) MarkSuccess() {
	cb.mu.Lock()
	cb.state = Closed
	cb.failures = 0
	cb.mu.Unlock()
}
func (cb *CircuitBreaker) MarkFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	if cb.failures >= cb.threshold {
		cb.state = Open
		cb.openUntil = time.Now().Add(cb.cooldown)
	}
}
