// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kopexa-grc/common/errors"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/rs/zerolog/log"
)

// EmbeddedServer is an embedded NATS server that implements the Shutdownable interface
type EmbeddedServer struct {
	server *server.Server
	mu     sync.Mutex

	ready chan struct{}
}

// NewEmbeddedServer creates a new embedded NATS server
func NewEmbeddedServer(opts *server.Options) (*EmbeddedServer, error) {
	if opts == nil {
		opts = &server.Options{
			Host:           "127.0.0.1",
			Port:           -1, // random port
			NoLog:          true,
			NoSigs:         true,
			MaxControlLine: 256,
		}
	}

	s, err := server.NewServer(opts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create NATS server")
	}

	return &EmbeddedServer{
		server: s,
	}, nil
}

// Start starts the NATS server
func (es *EmbeddedServer) Start() error {
	es.mu.Lock()
	defer es.mu.Unlock()

	if es.server == nil {
		return fmt.Errorf("server is nil")
	}

	es.ready = make(chan struct{})

	log.Info().Msg("EmbeddedServer: starting server goroutine")
	go func() {
		log.Info().Msg("EmbeddedServer: inside goroutine, calling Start()")
		es.server.Start()
	}()

	log.Info().Msg("EmbeddedServer: waiting for ReadyForConnections...")
	if !es.server.ReadyForConnections(4 * time.Second) {
		log.Error().Msg("EmbeddedServer: ReadyForConnections() timeout or failed")
		return fmt.Errorf("NATS server failed to start")
	}

	log.Info().Msg("EmbeddedServer: ReadyForConnections() succeeded")
	close(es.ready)

	return nil
}

// Ready returns a channel that is closed once the server and client are ready
func (es *EmbeddedServer) Ready() <-chan struct{} {
	es.mu.Lock()
	defer es.mu.Unlock()

	return es.ready
}

// Shutdown implements the Shutdownable interface
func (es *EmbeddedServer) Shutdown(ctx context.Context) error {
	es.mu.Lock()
	defer es.mu.Unlock()

	if es.server == nil {
		return nil
	}

	// stop the server
	es.server.Shutdown()

	// wait until the server is fully shutdown
	es.server.WaitForShutdown()

	return nil
}

// Addr returns the address of the server
func (es *EmbeddedServer) Addr() string {
	es.mu.Lock()
	defer es.mu.Unlock()

	if es.server == nil || !es.server.Running() {
		return ""
	}

	addr := es.server.Addr()
	if addr == nil {
		return ""
	}

	return addr.String()
}

// IsRunning returns true if the server is running
func (es *EmbeddedServer) IsRunning() bool {
	es.mu.Lock()
	defer es.mu.Unlock()

	return es.server != nil && es.server.Running()
}

// GetServer returns the internal NATS server
func (es *EmbeddedServer) GetServer() *server.Server {
	es.mu.Lock()
	defer es.mu.Unlock()

	return es.server
}
