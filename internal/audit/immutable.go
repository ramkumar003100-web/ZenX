package audit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)

type ImmutableStore struct {
	Store    Store
	lastHash string
}

func (s *ImmutableStore) Write(ctx context.Context, e Entry) error {
	b := []byte(e.Actor + e.Action + e.Resource + s.lastHash)
	h := sha256.Sum256(b)
	e.Metadata = map[string]any{"chain_hash": hex.EncodeToString(h[:])}
	s.lastHash = e.Metadata["chain_hash"].(string)
	return s.Store.Write(ctx, e)
}
func (s *ImmutableStore) Search(ctx context.Context, q string) ([]Entry, error) {
	return s.Store.Search(ctx, q)
}
