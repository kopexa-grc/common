// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package messaging

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a self-signed certificate for testing
func createTestCert(t *testing.T) (certPEM, keyPEM []byte, cleanup func()) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Test Organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	require.NoError(t, err)

	// Encode certificate and private key to PEM
	certPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return certPEM, keyPEM, func() {}
}

// Helper function to create temporary files for testing
func createTempFiles(t *testing.T, certPEM, keyPEM []byte) (certFile, keyFile string, cleanup func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "tls-test")
	require.NoError(t, err)

	// Create certificate file
	certFile = filepath.Join(tempDir, "cert.pem")
	err = os.WriteFile(certFile, certPEM, 0600)
	require.NoError(t, err)

	// Create key file
	keyFile = filepath.Join(tempDir, "key.pem")
	err = os.WriteFile(keyFile, keyPEM, 0600)
	require.NoError(t, err)

	return certFile, keyFile, func() {
		os.RemoveAll(tempDir)
	}
}

func TestTLSConfig_LoadCertFromValues(t *testing.T) {
	certPEM, keyPEM, cleanup := createTestCert(t)
	defer cleanup()

	cert, err := LoadCertFromValues(string(certPEM), string(keyPEM))
	require.NoError(t, err)
	assert.NotNil(t, cert)
	assert.Equal(t, 1, len(cert.Certificate))
}

func TestTLSConfig_LoadCertFromFiles(t *testing.T) {
	certPEM, keyPEM, cleanup := createTestCert(t)
	defer cleanup()

	certFile, keyFile, cleanupFiles := createTempFiles(t, certPEM, keyPEM)
	defer cleanupFiles()

	cert, err := LoadCertFromFiles(certFile, keyFile)
	require.NoError(t, err)
	assert.NotNil(t, cert)
	assert.Equal(t, 1, len(cert.Certificate))
}

func TestTLSConfig_LoadCAFromValue(t *testing.T) {
	certPEM, _, cleanup := createTestCert(t)
	defer cleanup()

	pool, err := LoadCAFromValue(string(certPEM))
	require.NoError(t, err)
	assert.NotNil(t, pool)
}

func TestTLSConfig_LoadCAFromFiles(t *testing.T) {
	certPEM, _, cleanup := createTestCert(t)
	defer cleanup()

	certFile, _, cleanupFiles := createTempFiles(t, certPEM, []byte{})
	defer cleanupFiles()

	pool, err := LoadCAFromFiles([]string{certFile})
	require.NoError(t, err)
	assert.NotNil(t, pool)
}

func TestTLSConfig_TLSConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      TLSConfig
		expectError bool
	}{
		{
			name: "Valid configuration with cert and key values",
			config: TLSConfig{
				Enabled: true,
				Cert:    "test-cert",
				Key:     "test-key",
			},
			expectError: true, // Will fail because cert/key are invalid
		},
		{
			name: "Valid configuration with cert and key files",
			config: TLSConfig{
				Enabled:  true,
				CertFile: "test-cert.pem",
				KeyFile:  "test-key.pem",
			},
			expectError: true, // Will fail because files don't exist
		},
		{
			name: "Disabled TLS",
			config: TLSConfig{
				Enabled: false,
			},
			expectError: false,
		},
		{
			name: "Insecure configuration",
			config: TLSConfig{
				Enabled:  true,
				Insecure: true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := tt.config.TLSConfig()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.config.Enabled {
					assert.NotNil(t, config)
				} else {
					assert.Nil(t, config)
				}
			}
		})
	}
}

func TestTLSConfig_TLSConfigWithValidCert(t *testing.T) {
	certPEM, keyPEM, cleanup := createTestCert(t)
	defer cleanup()

	config := TLSConfig{
		Enabled: true,
		Cert:    string(certPEM),
		Key:     string(keyPEM),
	}

	tlsConfig, err := config.TLSConfig()
	require.NoError(t, err)
	assert.NotNil(t, tlsConfig)
	assert.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Equal(t, 1, len(tlsConfig.Certificates))
}

func TestTLSConfig_TLSConfigWithValidFiles(t *testing.T) {
	certPEM, keyPEM, cleanup := createTestCert(t)
	defer cleanup()

	certFile, keyFile, cleanupFiles := createTempFiles(t, certPEM, keyPEM)
	defer cleanupFiles()

	config := TLSConfig{
		Enabled:  true,
		CertFile: certFile,
		KeyFile:  keyFile,
	}

	tlsConfig, err := config.TLSConfig()
	require.NoError(t, err)
	assert.NotNil(t, tlsConfig)
	assert.Equal(t, uint16(tls.VersionTLS12), tlsConfig.MinVersion)
	assert.False(t, tlsConfig.InsecureSkipVerify)
	assert.Equal(t, 1, len(tlsConfig.Certificates))
}
