package plugins

import (
	"context"
	"sync"
)

type Hook interface {
	OnRequest(context.Context)
	OnJobComplete(context.Context, string, error)
	OnShutdown(context.Context)
}

type Manager struct {
	mu      sync.RWMutex
	plugins []Hook
}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) Register(p Hook) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.plugins = append(m.plugins, p)
}

func (m *Manager) OnRequest(ctx context.Context) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.plugins {
		p.OnRequest(ctx)
	}
}

func (m *Manager) OnJobComplete(ctx context.Context, name string, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.plugins {
		p.OnJobComplete(ctx, name, err)
	}
}

func (m *Manager) OnShutdown(ctx context.Context) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.plugins {
		p.OnShutdown(ctx)
	}
}
