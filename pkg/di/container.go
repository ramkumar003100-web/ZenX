package di

import (
	"fmt"
	"reflect"
	"sync"
)

type Container struct {
	mu        sync.RWMutex
	singleton map[reflect.Type]any
}

func New() *Container {
	return &Container{singleton: make(map[reflect.Type]any)}
}

func (c *Container) Register(instance any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.singleton[reflect.TypeOf(instance)] = instance
}

func (c *Container) Resolve(target any) error {
	ptr := reflect.ValueOf(target)
	if ptr.Kind() != reflect.Ptr || ptr.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be pointer to struct")
	}

	st := ptr.Elem()
	c.mu.RLock()
	defer c.mu.RUnlock()

	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		if !f.CanSet() {
			continue
		}
		if dep, ok := c.singleton[f.Type()]; ok {
			f.Set(reflect.ValueOf(dep))
		}
	}
	return nil
}
