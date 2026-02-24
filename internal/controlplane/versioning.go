package controlplane

import "time"

type ConfigVersion struct {
	Version   int64             `json:"version"`
	CreatedAt time.Time         `json:"created_at"`
	Data      map[string]any    `json:"data"`
	Meta      map[string]string `json:"meta,omitempty"`
}
