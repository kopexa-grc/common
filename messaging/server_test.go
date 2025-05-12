// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedServer_New(t *testing.T) {
	tests := []struct {
		name    string
		opts    *server.Options
		wantErr bool
	}{
		{
			name:    "Mit Standard-Optionen",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "Mit benutzerdefinierten Optionen",
			opts: &server.Options{
				Host:           "127.0.0.1",
				Port:           -1,
				NoLog:          true,
				NoSigs:         true,
				MaxControlLine: 256,
				Username:       "testuser",
				Password:       "testpass",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es, err := NewEmbeddedServer(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, es)
			assert.NotNil(t, es.server)
		})
	}
}

func TestEmbeddedServer_Start(t *testing.T) {
	es, err := NewEmbeddedServer(nil)
	require.NoError(t, err)

	err = es.Start()
	require.NoError(t, err)

	assert.True(t, es.IsRunning())
	assert.NotEmpty(t, es.Addr())
}

func TestEmbeddedServer_ConcurrentAccess(t *testing.T) {
	es, err := NewEmbeddedServer(nil)
	require.NoError(t, err)

	// Starte den Server
	err = es.Start()
	require.NoError(t, err)

	// Führe mehrere gleichzeitige Operationen aus
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- struct{}{} }()

			// Versuche, eine Verbindung herzustellen
			nc, err := nats.Connect(es.Addr())
			if err != nil {
				return
			}
			defer nc.Close()

			// Führe einige NATS-Operationen aus
			sub, err := nc.SubscribeSync("test")
			if err != nil {
				return
			}

			err = nc.Publish("test", []byte("hello"))
			if err != nil {
				return
			}

			_, err = sub.NextMsg(time.Second)
			if err != nil {
				return
			}
		}()
	}

	// Warte auf alle Goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Führe den Shutdown durch
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = es.Shutdown(ctx)
	require.NoError(t, err)

	assert.False(t, es.IsRunning())
}
