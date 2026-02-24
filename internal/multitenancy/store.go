package multitenancy

import (
	"context"
	"sync"
)

type TenantConfig struct {
	Schema   string
	Features map[string]bool
	Metadata map[string]string
}

type ConfigStore struct {
	mu   sync.RWMutex
	data map[string]TenantConfig
}

func NewConfigStore() *ConfigStore { return &ConfigStore{data: map[string]TenantConfig{}} }
func (s *ConfigStore) Set(tenant string, c TenantConfig) {
	s.mu.Lock()
	s.data[tenant] = c
	s.mu.Unlock()
}
func (s *ConfigStore) Get(_ context.Context, tenant string) (TenantConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.data[tenant]
	return c, ok
}
