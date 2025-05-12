// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package messaging

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/kopexa-grc/common/errors"
)

var (
	ErrFailedToAddCACert     = errors.New(errors.BadRequest, "failed to add ca cert")
	ErrFailedToAddCACertFile = errors.New(errors.BadRequest, "failed to add ca cert at %s")
)

type TLSConfig struct {
	CAFiles  []string `mapstructure:"ca_files" envconfig:"ca_files" json:"ca_files" yaml:"ca_files"`
	KeyFile  string   `mapstructure:"key_file" split_words:"true" json:"key_file" yaml:"key_file"`
	CertFile string   `mapstructure:"cert_file" split_words:"true" json:"cert_file" yaml:"cert_file"`

	Cert string `mapstructure:"cert"`
	Key  string `mapstructure:"key"`
	CA   string `mapstructure:"ca"`

	Insecure bool `default:"false"`
	Enabled  bool `default:"false"`
}

func (cfg TLSConfig) TLSConfig() (*tls.Config, error) {
	var err error

	tlsConf := &tls.Config{
		MinVersion: tls.VersionTLS12,
		// nolint:gosec
		InsecureSkipVerify: cfg.Insecure,
	}

	// Load CA
	switch {
	case cfg.CA != "":
		tlsConf.RootCAs, err = LoadCAFromValue(cfg.CA)
	case len(cfg.CAFiles) > 0:
		tlsConf.RootCAs, err = LoadCAFromFiles(cfg.CAFiles)
	default:
		tlsConf.RootCAs, err = x509.SystemCertPool()
	}

	if err != nil {
		return nil, errors.Wrap(err, "error setting up root ca pool")
	}

	// Load Certs if any
	var cert tls.Certificate
	if cfg.Cert != "" && cfg.Key != "" {
		cert, err = LoadCertFromValues(cfg.Cert, cfg.Key)
		tlsConf.Certificates = append(tlsConf.Certificates, cert)
	} else if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err = LoadCertFromFiles(cfg.CertFile, cfg.KeyFile)
		tlsConf.Certificates = append(tlsConf.Certificates, cert)
	}

	if err != nil {
		return nil, errors.Wrap(err, "error loading certificate keypair")
	}

	// Backwards compatibility: if TLS is not explicitly enabled, return nil if no certificate was provided
	// Old code disabled TLS by not providing a certificate, which returned nil when calling TLSConfig()
	if !cfg.Enabled && len(tlsConf.Certificates) == 0 {
		return nil, nil
	}

	return tlsConf, nil
}

func LoadCertFromValues(certPEM, keyPEM string) (tls.Certificate, error) {
	return tls.X509KeyPair([]byte(certPEM), []byte(keyPEM))
}

func LoadCertFromFiles(certFile, keyFile string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFile, keyFile)
}

func LoadCAFromFiles(cafiles []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()

	for _, caFile := range cafiles {
		caData, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		if !pool.AppendCertsFromPEM(caData) {
			return nil, fmt.Errorf("%w: %s", ErrFailedToAddCACertFile, caFile)
		}
	}

	return pool, nil
}

func LoadCAFromValue(ca string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(ca)) {
		return nil, ErrFailedToAddCACert
	}

	return pool, nil
}
