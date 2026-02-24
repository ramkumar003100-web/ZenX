package cluster

import (
	"context"
	"sort"
	"sync"
	"time"
)

type LeaderElector struct {
	mu       sync.RWMutex
	leaderID string
	members  *Membership
}

func NewLeaderElector(m *Membership) *LeaderElector { return &LeaderElector{members: m} }

func (e *LeaderElector) Run(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			e.elect()
		}
	}
}

func (e *LeaderElector) Leader() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.leaderID
}

func (e *LeaderElector) elect() {
	nodes := e.members.Nodes()
	healthy := make([]string, 0)
	for _, n := range nodes {
		if n.IsHealthy {
			healthy = append(healthy, n.ID)
		}
	}
	sort.Strings(healthy)
	e.mu.Lock()
	defer e.mu.Unlock()
	if len(healthy) == 0 {
		e.leaderID = ""
		return
	}
	e.leaderID = healthy[0]
}
