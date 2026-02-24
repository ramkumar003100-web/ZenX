package oauth

import "net/http"

func OIDCDiscovery(issuer string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_ = jsonWrite(w, map[string]any{"issuer": issuer, "authorization_endpoint": issuer + "/authorize", "token_endpoint": issuer + "/token", "jwks_uri": issuer + "/jwks"})
	}
}
