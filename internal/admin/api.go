package admin

import (
	"encoding/json"
	"net/http"

	"zenx/internal/featureflags"
)

type API struct {
	FlagManager *featureflags.Manager
	QueueDepth  func() map[string]int
	Tenants     func() []string
	WSStats     func() map[string]int
}

func (a API) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/admin/system/metrics", a.systemMetrics)
	mux.HandleFunc("/admin/feature-flags", a.featureFlags)
	mux.HandleFunc("/admin/tenants", a.tenants)
	mux.HandleFunc("/admin/jobs", a.jobs)
	mux.HandleFunc("/admin/websocket", a.websocket)
}

func (a API) systemMetrics(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
}
func (a API) featureFlags(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]any{"supported": true})
}
func (a API) tenants(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(a.Tenants())
}
func (a API) jobs(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(a.QueueDepth())
}
func (a API) websocket(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(a.WSStats())
}
