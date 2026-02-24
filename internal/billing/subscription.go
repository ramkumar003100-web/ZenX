package billing

import (
	"context"
	"sync"
)

type Subscription struct {
	TenantID string
	Plan     string
	Active   bool
}

type Service struct {
	mu   sync.RWMutex
	subs map[string]Subscription
}

func NewService() *Service { return &Service{subs: map[string]Subscription{}} }
func (s *Service) Set(_ context.Context, sub Subscription) {
	s.mu.Lock()
	s.subs[sub.TenantID] = sub
	s.mu.Unlock()
}
func (s *Service) Get(_ context.Context, tenant string) (Subscription, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.subs[tenant]
	return v, ok
}
