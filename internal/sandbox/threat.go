package sandbox

import (
	"sync"
	"time"
)

type ThreatScore struct {
	IP          string
	Score       float64
	UpdatedAt   time.Time
	ReasonCodes []string
}

type ReputationCache struct {
	mu    sync.RWMutex
	items map[string]ThreatScore
}

func NewReputationCache() *ReputationCache   { return &ReputationCache{items: map[string]ThreatScore{}} }
func (c *ReputationCache) Set(t ThreatScore) { c.mu.Lock(); c.items[t.IP] = t; c.mu.Unlock() }
func (c *ReputationCache) Get(ip string) (ThreatScore, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	t, ok := c.items[ip]
	return t, ok
}
