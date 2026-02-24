package sandbox

import "sync"

type Policy struct {
	AllowDB        bool
	AllowNetwork   bool
	AllowFileRead  bool
	AllowFileWrite bool
	MaxCPUPercent  float64
}

type EnginePolicyStore struct {
	mu   sync.RWMutex
	data map[string]Policy
}

func NewPolicyStore() *EnginePolicyStore { return &EnginePolicyStore{data: map[string]Policy{}} }
func (s *EnginePolicyStore) Set(module string, p Policy) {
	s.mu.Lock()
	s.data[module] = p
	s.mu.Unlock()
}
func (s *EnginePolicyStore) Get(module string) (Policy, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.data[module]
	return p, ok
}
