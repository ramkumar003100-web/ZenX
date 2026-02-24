package oauth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type OAuth2Server struct {
	mu      sync.Mutex
	clients map[string]string
	tokens  map[string]time.Time
}

func NewOAuth2Server() *OAuth2Server {
	return &OAuth2Server{clients: map[string]string{}, tokens: map[string]time.Time{}}
}
func (s *OAuth2Server) RegisterClient(id, secret string) {
	s.mu.Lock()
	s.clients[id] = secret
	s.mu.Unlock()
}
func (s *OAuth2Server) TokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, sec, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "missing basic auth", http.StatusUnauthorized)
			return
		}
		s.mu.Lock()
		valid := s.clients[id] == sec
		s.mu.Unlock()
		if !valid {
			http.Error(w, "invalid client", http.StatusUnauthorized)
			return
		}
		b := make([]byte, 32)
		_, _ = rand.Read(b)
		t := hex.EncodeToString(b)
		s.mu.Lock()
		s.tokens[t] = time.Now().Add(time.Hour)
		s.mu.Unlock()
		_ = jsonWrite(w, Token{AccessToken: t, ExpiresIn: 3600, TokenType: "Bearer"})
	}
}
func (s *OAuth2Server) Introspect(token string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.tokens[token]
	return ok && time.Now().Before(exp)
}
func jsonWrite(w http.ResponseWriter, v any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
