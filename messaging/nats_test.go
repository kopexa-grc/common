// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package messaging

import (
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats-server/v2/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create and start an embedded NATS server
func runNATSServer(_ *testing.T, opts *server.Options) (*server.Server, error) {
	if opts == nil {
		opts = &server.Options{
			Host:           "127.0.0.1",
			Port:           -1, // Random port
			NoLog:          true,
			NoSigs:         true,
			MaxControlLine: 256,
		}
	}

	s := test.RunServer(opts)

	return s, nil
}

func TestNATSClient_NoAuth(t *testing.T) {
	// Start NATS server
	s, err := runNATSServer(t, nil)
	require.NoError(t, err)
	defer s.Shutdown()

	// Create client config
	cfg := &NATSConfig{
		Servers: []string{s.Addr().String()},
		TLS:     TLSConfig{Enabled: false},
	}

	// Connect
	nc, err := NewNATSClient(cfg)
	require.NoError(t, err)
	defer nc.Close()

	assert.True(t, nc.IsConnected())
}

func TestNATSClient_UserAuth(t *testing.T) {
	// Start NATS server with user authentication
	opts := &server.Options{
		Host:           "127.0.0.1",
		Port:           -1,
		NoLog:          true,
		NoSigs:         true,
		MaxControlLine: 256,
		Username:       "testuser",
		Password:       "testpass",
	}

	s, err := runNATSServer(t, opts)
	require.NoError(t, err)
	defer s.Shutdown()

	tests := []struct {
		name        string
		config      *NATSConfig
		expectError bool
	}{
		{
			name: "Valid credentials",
			config: &NATSConfig{
				Servers: []string{s.Addr().String()},
				Auth: NatsAuth{
					Method:   NatsAuthMethodUser,
					User:     "testuser",
					Password: "testpass",
				},
			},
			expectError: false,
		},
		{
			name: "Invalid credentials",
			config: &NATSConfig{
				Servers: []string{s.Addr().String()},
				Auth: NatsAuth{
					Method:   NatsAuthMethodUser,
					User:     "wronguser",
					Password: "wrongpass",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc, err := NewNATSClient(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			defer nc.Close()
			assert.True(t, nc.IsConnected())
		})
	}
}

func TestNATSClient_TokenAuth(t *testing.T) {
	// Start NATS server with token authentication
	opts := &server.Options{
		Host:           "127.0.0.1",
		Port:           -1,
		NoLog:          true,
		NoSigs:         true,
		MaxControlLine: 256,
		Authorization:  "testtoken",
	}

	s, err := runNATSServer(t, opts)
	require.NoError(t, err)
	defer s.Shutdown()

	tests := []struct {
		name        string
		config      *NATSConfig
		expectError bool
	}{
		{
			name: "Valid token",
			config: &NATSConfig{
				Servers: []string{s.Addr().String()},
				Auth: NatsAuth{
					Method: NatsAuthMethodToken,
					Token:  "testtoken",
				},
			},
			expectError: false,
		},
		{
			name: "Invalid token",
			config: &NATSConfig{
				Servers: []string{s.Addr().String()},
				Auth: NatsAuth{
					Method: NatsAuthMethodToken,
					Token:  "wrongtoken",
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nc, err := NewNATSClient(tt.config)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			defer nc.Close()
			assert.True(t, nc.IsConnected())
		})
	}
}
