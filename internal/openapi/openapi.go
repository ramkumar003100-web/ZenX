package openapi

import "encoding/json"

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
	if b.spec.Paths[path] == nil {
		b.spec.Paths[path] = map[string]Path{}
	}
	p := Path{Summary: summary, Responses: map[string]any{"200": map[string]string{"description": "OK"}}}
	if secured {
		p.Security = []map[string][]string{{"BearerAuth": roles}}
	}
	b.spec.Paths[path][method] = p
}

func (b *Builder) JSON() ([]byte, error) {
	return json.MarshalIndent(b.spec, "", "  ")
}
