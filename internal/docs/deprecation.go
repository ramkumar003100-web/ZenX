package docs

import "net/http"

func DeprecationHeaders(version, sunset string) http.Header {
	h := http.Header{}
	h.Set("X-API-Version", version)
	h.Set("Deprecation", "true")
	h.Set("Sunset", sunset)
	return h
}
