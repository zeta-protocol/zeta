package gateway

import (
	"crypto/tls"
	"fmt"

	"github.com/zeta-protocol/zeta/paths"
	"golang.org/x/crypto/acme/autocert"
)

func GenerateTlsConfig(g *Config, zetaPaths paths.Paths) (*tls.Config, error) {
	if g.HTTPSEnabled {
		if g.AutoCertDomain != "" {
			if g.CertificateFile != "" || g.KeyFile != "" {
				return nil, fmt.Errorf("Autocert is enabled, and a pre-generated certificate/key specified; use one or the other")
			}
			certDir := zetaPaths.StatePathFor(paths.DataNodeAutoCertHome)

			certManager := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(g.AutoCertDomain),
				Cache:      autocert.DirCache(certDir),
			}
			return &tls.Config{
				GetCertificate: certManager.GetCertificate,
				NextProtos:     []string{"http/1.1", "acme-tls/1"},
			}, nil
		}

		certificate, err := tls.LoadX509KeyPair(g.CertificateFile, g.KeyFile)
		if err != nil {
			return nil, err
		}
		certificates := []tls.Certificate{certificate}
		return &tls.Config{
			Certificates: certificates,
		}, nil
	}

	return nil, nil
}
