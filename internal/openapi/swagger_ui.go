package openapi

import (
	"fmt"
	"net/http"
)

func SwaggerUI(specPath string) http.HandlerFunc {
	html := fmt.Sprintf(`<!doctype html>
<html>
<head><title>ZenX Swagger</title><link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" /></head>
<body><div id="swagger-ui"></div><script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>window.ui = SwaggerUIBundle({url: '%s', dom_id: '#swagger-ui'})</script></body>
</html>`, specPath)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(html))
	}
}
