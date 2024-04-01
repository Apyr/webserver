package tlsutils

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/idna"
)

func NewTLSConfig(cache autocert.Cache, handler http.Handler, acmeHosts, selfHosts []string) (*tls.Config, http.Handler) {
	selfSignedHosts := make(map[string]struct{}, len(selfHosts))
	for _, host := range selfHosts {
		name, err := idna.Lookup.ToASCII(host)
		if err != nil {
			name = host
		}
		selfSignedHosts[name] = struct{}{}
	}

	manager := autocert.Manager{
		Cache:      cache,
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(acmeHosts...),
	}
	handler = manager.HTTPHandler(handler)
	config := manager.TLSConfig()

	config.GetCertificate = func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cert, err := manager.GetCertificate(chi)
		if err != nil {
			name, err := idna.Lookup.ToASCII(chi.ServerName)
			if err != nil {
				return nil, err
			}
			if _, ok := selfSignedHosts[name]; !ok {
				return nil, err
			}
			return generateCert(cache, name)
		}
		return cert, nil
	}

	return config, handler
}
