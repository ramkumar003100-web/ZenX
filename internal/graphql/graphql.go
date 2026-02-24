package graphql

import (
	"context"
	"encoding/json"
	"net/http"
)

type Resolver func(context.Context, map[string]any) (any, error)

type Engine struct {
	queries   map[string]Resolver
	mutations map[string]Resolver
	authorize func(context.Context, string) bool
}

func New(authorize func(context.Context, string) bool) *Engine {
	return &Engine{queries: map[string]Resolver{}, mutations: map[string]Resolver{}, authorize: authorize}
}

func (e *Engine) RegisterQuery(name string, r Resolver)    { e.queries[name] = r }
func (e *Engine) RegisterMutation(name string, r Resolver) { e.mutations[name] = r }

func (e *Engine) Handler() http.HandlerFunc {
	type req struct {
		Type string         `json:"type"`
		Name string         `json:"name"`
		Args map[string]any `json:"args"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var in req
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if !e.authorize(r.Context(), in.Name) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		var fn Resolver
		if in.Type == "mutation" {
			fn = e.mutations[in.Name]
		} else {
			fn = e.queries[in.Name]
		}
		if fn == nil {
			http.Error(w, "operation not found", http.StatusNotFound)
			return
		}
		result, err := fn(r.Context(), in.Args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": result})
	}
}
