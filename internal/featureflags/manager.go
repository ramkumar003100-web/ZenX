package featureflags

import "sync"

type FlagRule struct {
	Global bool
	Roles  map[string]bool
	Users  map[string]bool
}

type Manager struct {
	mu    sync.RWMutex
	flags map[string]FlagRule
}

func NewManager() *Manager {
	return &Manager{flags: map[string]FlagRule{}}
}

func (m *Manager) Set(flag string, rule FlagRule) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if rule.Roles == nil {
		rule.Roles = map[string]bool{}
	}
	if rule.Users == nil {
		rule.Users = map[string]bool{}
	}
	m.flags[flag] = rule
}

func (m *Manager) Enabled(flag, userID string, roles ...string) bool {
	m.mu.RLock()
	rule, ok := m.flags[flag]
	m.mu.RUnlock()
	if !ok {
		return false
	}
	if rule.Global || rule.Users[userID] {
		return true
	}
	for _, role := range roles {
		if rule.Roles[role] {
			return true
		}
	}
	return false
}
