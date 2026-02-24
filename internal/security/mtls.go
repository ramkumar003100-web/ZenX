package security

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

func NewMTLSServerConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	caPem, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caPem)
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{ClientCAs: pool, ClientAuth: tls.RequireAndVerifyClientCert, Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13}, nil
}
