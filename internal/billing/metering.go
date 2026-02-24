package billing

import "sync"

type Meter struct {
	mu    sync.Mutex
	usage map[string]int64
}

func NewMeter() *Meter                          { return &Meter{usage: map[string]int64{}} }
func (m *Meter) Add(tenant string, units int64) { m.mu.Lock(); m.usage[tenant] += units; m.mu.Unlock() }
func (m *Meter) Get(tenant string) int64        { m.mu.Lock(); defer m.mu.Unlock(); return m.usage[tenant] }
