package cluster

import (
	"context"
	"sync"
	"time"
)

type Node struct {
	ID        string
	Address   string
	Region    string
	LastBeat  time.Time
	IsHealthy bool
}

type Membership struct {
	mu    sync.RWMutex
	nodes map[string]Node
	ttl   time.Duration
}

func NewMembership(ttl time.Duration) *Membership {
	return &Membership{nodes: map[string]Node{}, ttl: ttl}
}

func (m *Membership) Heartbeat(_ context.Context, n Node) {
	m.mu.Lock()
	defer m.mu.Unlock()
	n.LastBeat = time.Now()
	n.IsHealthy = true
	m.nodes[n.ID] = n
}

func (m *Membership) Nodes() []Node {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	out := make([]Node, 0, len(m.nodes))
	for id, n := range m.nodes {
		if now.Sub(n.LastBeat) > m.ttl {
			n.IsHealthy = false
			m.nodes[id] = n
		}
		out = append(out, n)
	}
	return out
}
