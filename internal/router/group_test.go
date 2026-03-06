package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupVersionAndMiddleware(t *testing.T) {
	r := New()

	api := r.Group("/api", func(next HandlerFunc) HandlerFunc {
		return func(c *Context) {
			c.Writer.Header().Set("X-Group", "applied")
			next(c)
		}
	})

	api.Version("1").Handle(http.MethodGet, "/users/:id", func(c *Context) {
		if c.Param("id") == "" {
			t.Fatalf("expected id param")
		}
		c.Writer.WriteHeader(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/42", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", rec.Code)
	}
	if got := rec.Header().Get("X-Group"); got != "applied" {
		t.Fatalf("expected group middleware header, got %q", got)
	}
}
