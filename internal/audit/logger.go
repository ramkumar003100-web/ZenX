package audit

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Entry struct {
	Time     time.Time      `json:"time"`
	Actor    string         `json:"actor"`
	Action   string         `json:"action"`
	Resource string         `json:"resource"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

type Store interface {
	Write(context.Context, Entry) error
	Search(context.Context, string) ([]Entry, error)
}

type FileStore struct {
	mu   sync.Mutex
	path string
}

func NewFileStore(path string) *FileStore { return &FileStore{path: path} }
func (s *FileStore) Write(_ context.Context, e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(e)
}
func (s *FileStore) Search(_ context.Context, _ string) ([]Entry, error) { return nil, nil }
