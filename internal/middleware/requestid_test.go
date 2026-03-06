package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"zenx/internal/router"
)

func TestRequestIDGeneratesAndPropagates(t *testing.T) {
	r := router.New()
	r.Use(RequestID(""))

	r.Handle(http.MethodGet, "/", func(c *router.Context) {
		if rid := GetRequestID(c.Request.Context()); rid == "" {
			t.Fatalf("expected request id in context")
		}
		c.Writer.WriteHeader(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", rec.Code)
	}
	if got := rec.Header().Get("X-Request-ID"); got == "" {
		t.Fatalf("expected request id response header")
	}
}

func TestRequestIDUsesIncomingHeader(t *testing.T) {
	r := router.New()
	r.Use(RequestID(""))

	r.Handle(http.MethodGet, "/", func(c *router.Context) {
		c.Writer.WriteHeader(http.StatusNoContent)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "abc-123")
	r.ServeHTTP(rec, req)

	if got := rec.Header().Get("X-Request-ID"); got != "abc-123" {
		t.Fatalf("expected incoming request id to be preserved, got %q", got)
	}
}
