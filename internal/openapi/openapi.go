package openapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

type Spec struct {
	OpenAPI    string                     `json:"openapi"`
	Info       map[string]string          `json:"info"`
	Paths      map[string]map[string]Path `json:"paths"`
	Components map[string]any             `json:"components,omitempty"`
}

type Path struct {
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	Tags        []string              `json:"tags,omitempty"`
	Security    []map[string][]string `json:"security,omitempty"`
	Responses   map[string]any        `json:"responses"`
}

type Builder struct {
	mu   sync.RWMutex
	spec Spec
}

func New(title, version string) *Builder {
	return &Builder{spec: Spec{
		OpenAPI: "3.0.3",
		Info:    map[string]string{"title": title, "version": version},
		Paths:   map[string]map[string]Path{},
		Components: map[string]any{
			"securitySchemes": map[string]any{
				"BearerAuth": map[string]string{"type": "http", "scheme": "bearer", "bearerFormat": "JWT"},
			},
		},
	}}
}

func (b *Builder) AddPath(path, method, summary string, secured bool, roles ...string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.spec.Paths[path] == nil {
		b.spec.Paths[path] = map[string]Path{}
	}
	method = strings.ToLower(method)
	p := Path{Summary: summary, Responses: map[string]any{"200": map[string]string{"description": "OK"}}}
	if secured {
		p.Security = []map[string][]string{{"BearerAuth": roles}}
	}
	b.spec.Paths[path][method] = p
}

func (b *Builder) JSON() ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return json.MarshalIndent(b.spec, "", "  ")
}

func (b *Builder) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		payload, err := b.JSON()
		if err != nil {
			http.Error(w, "spec build failed", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(payload)
	}
}