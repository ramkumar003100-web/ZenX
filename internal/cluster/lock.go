package cluster

import (
	"context"
	"sync"
	"time"
)

type DistributedLock struct {
	mu    sync.Mutex
	locks map[string]lockEntry
}

type lockEntry struct {
	owner string
	exp   time.Time
}

func NewDistributedLock() *DistributedLock { return &DistributedLock{locks: map[string]lockEntry{}} }

func (l *DistributedLock) Acquire(_ context.Context, key, owner string, ttl time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if e, ok := l.locks[key]; ok && now.Before(e.exp) {
		return false
	}
	l.locks[key] = lockEntry{owner: owner, exp: now.Add(ttl)}
	return true
}

func (l *DistributedLock) Release(_ context.Context, key, owner string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	e, ok := l.locks[key]
	if !ok || e.owner != owner {
		return false
	}
	delete(l.locks, key)
	return true
}
