package traffic

import (
	"context"
	"sort"
	"sync"
)

type Node struct {
	ID          string
	Region      string
	LatencyMS   float64
	LoadPercent float64
	Draining    bool
}

type Balancer struct {
	mu    sync.RWMutex
	nodes map[string]Node
}

func NewBalancer() *Balancer      { return &Balancer{nodes: map[string]Node{}} }
func (b *Balancer) Upsert(n Node) { b.mu.Lock(); b.nodes[n.ID] = n; b.mu.Unlock() }
func (b *Balancer) Drain(id string, draining bool) {
	b.mu.Lock()
	n := b.nodes[id]
	n.Draining = draining
	b.nodes[id] = n
	b.mu.Unlock()
}

func (b *Balancer) Pick(_ context.Context, region string) (Node, bool) {
	b.mu.RLock()
	list := make([]Node, 0, len(b.nodes))
	for _, n := range b.nodes {
		if !n.Draining {
			if region == "" || n.Region == region {
				list = append(list, n)
			}
		}
	}
	b.mu.RUnlock()
	if len(list) == 0 {
		return Node{}, false
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].LoadPercent == list[j].LoadPercent {
			return list[i].LatencyMS < list[j].LatencyMS
		}
		return list[i].LoadPercent < list[j].LoadPercent
	})
	return list[0], true
}
