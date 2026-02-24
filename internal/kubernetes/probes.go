package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
)

type Probe struct{ ready atomic.Bool }

func NewProbe() *Probe           { p := &Probe{}; p.ready.Store(true); return p }
func (p *Probe) SetReady(v bool) { p.ready.Store(v) }
func (p *Probe) Readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if !p.ready.Load() {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	}
}
func (p *Probe) Liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("alive"))
	}
}
func ConfigMapSecretHandler(configProvider func(context.Context) map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(configProvider(r.Context()))
	}
}
