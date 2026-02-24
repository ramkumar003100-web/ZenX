package runtime

import (
	"context"
	"sync"
	"time"
)

type Subsystem interface {
	Name() string
	Start(context.Context) error
	Health(context.Context) (float64, error)
}

type Supervisor struct {
	mu         sync.RWMutex
	subsystems map[string]Subsystem
	states     map[string]float64
	recovery   RecoveryPolicy
}

func NewSupervisor(recovery RecoveryPolicy) *Supervisor {
	return &Supervisor{subsystems: map[string]Subsystem{}, states: map[string]float64{}, recovery: recovery}
}

func (s *Supervisor) Register(ss Subsystem) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subsystems[ss.Name()] = ss
}

func (s *Supervisor) Run(ctx context.Context) {
	s.mu.RLock()
	list := make([]Subsystem, 0, len(s.subsystems))
	for _, ss := range s.subsystems {
		list = append(list, ss)
	}
	s.mu.RUnlock()

	for _, ss := range list {
		go s.runIsolated(ctx, ss)
	}

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			s.refreshHealth(ctx)
		}
	}
}

func (s *Supervisor) HealthSnapshot() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]float64, len(s.states))
	for k, v := range s.states {
		out[k] = v
	}
	return out
}

func (s *Supervisor) runIsolated(ctx context.Context, ss Subsystem) {
	for {
		if ctx.Err() != nil {
			return
		}
		err := RunWithRecovery(ctx, ss.Name(), func(runCtx context.Context) error {
			return ss.Start(runCtx)
		})
		if err == nil || !s.recovery.ShouldRestart(ss.Name(), err) {
			return
		}
		time.Sleep(s.recovery.Backoff(ss.Name()))
	}
}

func (s *Supervisor) refreshHealth(ctx context.Context) {
	s.mu.RLock()
	list := make([]Subsystem, 0, len(s.subsystems))
	for _, ss := range s.subsystems {
		list = append(list, ss)
	}
	s.mu.RUnlock()

	snap := map[string]float64{}
	for _, ss := range list {
		h, err := ss.Health(ctx)
		if err != nil {
			h = 0
		}
		snap[ss.Name()] = h
	}

	s.mu.Lock()
	s.states = snap
	s.mu.Unlock()
}
