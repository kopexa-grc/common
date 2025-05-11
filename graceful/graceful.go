// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package graceful

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// Shutdownable is a target that can be closed gracefully
type Shutdownable interface {
	Shutdown(context.Context) error
}

type target struct {
	name    string
	shut    Shutdownable
	timeout time.Duration
}

// Closer handles shutdown of servers and connections
type Closer struct {
	targets      []target
	targetsMutex sync.Mutex

	done     chan struct{}
	doneBool int32

	// for tests: if true, os.Exit is not called
	skipExit bool
}

// NewCloser creates a new Closer
func NewCloser() *Closer {
	return &Closer{
		done: make(chan struct{}),
	}
}

// Register inserts a target to shutdown gracefully
func (cc *Closer) Register(name string, shut Shutdownable, timeout time.Duration) {
	cc.targetsMutex.Lock()
	cc.targets = append(cc.targets, target{
		name:    name,
		shut:    shut,
		timeout: timeout,
	})
	cc.targetsMutex.Unlock()
}

// GetTargets returns a copy of the registered targets (for testing)
func (cc *Closer) GetTargets() []target {
	cc.targetsMutex.Lock()
	defer cc.targetsMutex.Unlock()

	targets := make([]target, len(cc.targets))
	copy(targets, cc.targets)
	return targets
}

// DetectShutdown asynchronously waits for a shutdown signal and then shuts down gracefully
// Returns a function to trigger a shutdown from the outside and a ready channel that signals when the detector is ready
func (cc *Closer) DetectShutdown() (trigger func(), ready chan struct{}) {
	readyChan := make(chan struct{})

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		defer signal.Stop(signals)

		// Signal that we're ready to handle signals
		close(readyChan)

		select {
		case sig := <-signals:
			log.Info().Str("signal", sig.String()).Msg("Triggering shutdown from signal")
		case <-cc.done:
			log.Info().Msg("Shutting down...")
		}

		if atomic.LoadInt32(&cc.doneBool) == 1 {
			return
		}

		if atomic.SwapInt32(&cc.doneBool, 1) != 1 {
			wg := sync.WaitGroup{}
			cc.targetsMutex.Lock()
			for _, targ := range cc.targets {
				wg.Add(1)
				go func(targ target) {
					defer wg.Done()

					ctx, cancel := context.WithTimeout(context.Background(), targ.timeout)
					defer cancel()

					if err := targ.shut.Shutdown(ctx); err != nil {
						if ctx.Err() == context.DeadlineExceeded {
							log.Warn().
								Str("target", targ.name).
								Dur("timeout", targ.timeout).
								Msg("Shutdown timed out")
						} else {
							log.Error().
								Err(err).
								Str("target", targ.name).
								Msg("Shutdown failed")
						}
					} else {
						log.Info().Str("target", targ.name).Msg("Shutdown finished")
					}
				}(targ)
			}
			cc.targetsMutex.Unlock()
			wg.Wait()

			if !cc.skipExit {
				os.Exit(0)
			}
		}
	}()

	return func() {
		cc.done <- struct{}{}
	}, readyChan
}

// SetSkipExit sets the skipExit option for tests
func (cc *Closer) SetSkipExit(skip bool) {
	cc.skipExit = skip
}

// ShutdownableFunc adapts a function into a Shutdownable.
type ShutdownableFunc func(context.Context) error

func (f ShutdownableFunc) Shutdown(ctx context.Context) error {
	return f(ctx)
}

// HTTPServerShutdown wraps an *http.Server into a graceful.Shutdownable
func HTTPServerShutdown(srv *http.Server) Shutdownable {
	return ShutdownableFunc(func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	})
}
