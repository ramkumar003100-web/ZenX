package controlplane

import (
	"context"
	"sync"
	"time"
)

type Controller struct {
	mu       sync.RWMutex
	versions []ConfigVersion
	current  ConfigVersion
}

func NewController() *Controller {
	v := ConfigVersion{Version: 1, CreatedAt: time.Now(), Data: map[string]any{}}
	return &Controller{versions: []ConfigVersion{v}, current: v}
}

func (c *Controller) Update(_ context.Context, data map[string]any, meta map[string]string) ConfigVersion {
	c.mu.Lock()
	defer c.mu.Unlock()
	v := ConfigVersion{Version: c.current.Version + 1, CreatedAt: time.Now(), Data: clone(data), Meta: cloneMeta(meta)}
	c.current = v
	c.versions = append(c.versions, v)
	return v
}

func (c *Controller) Rollback(_ context.Context, version int64) (ConfigVersion, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, v := range c.versions {
		if v.Version == version {
			c.current = v
			return v, true
		}
	}
	return ConfigVersion{}, false
}

func (c *Controller) Current() ConfigVersion {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.current
}

func clone(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
func cloneMeta(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
