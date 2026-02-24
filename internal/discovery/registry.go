package discovery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Tags     []string          `json:"tags,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Healthy  bool              `json:"healthy"`
}

type Registry interface {
	Register(context.Context, ServiceInstance) error
	Deregister(context.Context, string, string) error
	Instances(context.Context, string) ([]ServiceInstance, error)
}

type InMemoryRegistry struct {
	mu        sync.RWMutex
	instances map[string]map[string]ServiceInstance
}

func NewInMemoryRegistry() *InMemoryRegistry {
	return &InMemoryRegistry{instances: map[string]map[string]ServiceInstance{}}
}

func (r *InMemoryRegistry) Register(_ context.Context, inst ServiceInstance) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.instances[inst.Name] == nil {
		r.instances[inst.Name] = map[string]ServiceInstance{}
	}
	inst.Healthy = true
	r.instances[inst.Name][inst.ID] = inst
	return nil
}

func (r *InMemoryRegistry) Deregister(_ context.Context, name, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.instances[name] != nil {
		delete(r.instances[name], id)
	}
	return nil
}

func (r *InMemoryRegistry) Instances(_ context.Context, name string) ([]ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]ServiceInstance, 0, len(r.instances[name]))
	for _, i := range r.instances[name] {
		if i.Healthy {
			res = append(res, i)
		}
	}
	return res, nil
}

type ClientSideLB struct{ Registry Registry }

func (lb ClientSideLB) Pick(ctx context.Context, service string) (ServiceInstance, error) {
	instances, err := lb.Registry.Instances(ctx, service)
	if err != nil || len(instances) == 0 {
		return ServiceInstance{}, errors.New("no service instances")
	}
	return instances[rand.Intn(len(instances))], nil
}

type HTTPRegistry struct {
	BaseURL string
	Client  *http.Client
}

func (r *HTTPRegistry) Register(ctx context.Context, inst ServiceInstance) error {
	return r.do(ctx, http.MethodPut, fmt.Sprintf("%s/services/%s/%s", r.BaseURL, inst.Name, inst.ID), inst)
}
func (r *HTTPRegistry) Deregister(ctx context.Context, name, id string) error {
	return r.do(ctx, http.MethodDelete, fmt.Sprintf("%s/services/%s/%s", r.BaseURL, name, id), nil)
}
func (r *HTTPRegistry) Instances(ctx context.Context, name string) ([]ServiceInstance, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/services/%s", r.BaseURL, name), nil)
	resp, err := r.http().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out []ServiceInstance
	return out, json.NewDecoder(resp.Body).Decode(&out)
}
func (r *HTTPRegistry) do(ctx context.Context, method, url string, body any) error {
	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req, _ = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequestWithContext(ctx, method, url, nil)
	}
	resp, err := r.http().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("registry returned %d", resp.StatusCode)
	}
	return nil
}
func (r *HTTPRegistry) http() *http.Client {
	if r.Client != nil {
		return r.Client
	}
	return &http.Client{Timeout: 5 * time.Second}
}
