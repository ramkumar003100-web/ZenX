package performance

import (
	"crypto/tls"
	"net/http"
)

func HTTP2Server(addr string, h http.Handler) *http.Server {
	return &http.Server{Addr: addr, Handler: h, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
}
