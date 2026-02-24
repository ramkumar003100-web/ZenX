package plugins

import "sync"

type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Descriptor
}

func NewRegistry() *Registry { return &Registry{plugins: map[string]Descriptor{}} }
func (r *Registry) Register(d Descriptor) error {
	if err := d.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	r.plugins[d.Name] = d
	r.mu.Unlock()
	return nil
}
func (r *Registry) Get(name string) (Descriptor, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.plugins[name]
	return d, ok
}
