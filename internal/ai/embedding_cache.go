package ai

import "sync"

type EmbeddingCache struct {
	mu sync.RWMutex
	m  map[string][]float32
}

func NewEmbeddingCache() *EmbeddingCache { return &EmbeddingCache{m: map[string][]float32{}} }
func (c *EmbeddingCache) Get(k string) ([]float32, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[k]
	return v, ok
}
func (c *EmbeddingCache) Set(k string, v []float32) { c.mu.Lock(); c.m[k] = v; c.mu.Unlock() }
