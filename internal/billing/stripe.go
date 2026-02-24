package billing

import (
	"context"
	"encoding/json"
	"net/http"
)

type StripeEvent struct {
	ID   string         `json:"id"`
	Type string         `json:"type"`
	Data map[string]any `json:"data"`
}

func StripeWebhook(handler func(context.Context, StripeEvent) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var e StripeEvent
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
			http.Error(w, "bad payload", http.StatusBadRequest)
			return
		}
		if err := handler(r.Context(), e); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
