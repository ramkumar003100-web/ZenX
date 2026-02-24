package oauth

import "net/http"

func SSOLoginRedirect(providerURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, providerURL, http.StatusFound) }
}
