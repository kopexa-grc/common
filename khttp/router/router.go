// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/kopexa-grc/common/khttp/metric"
	"github.com/kopexa-grc/common/khttp/pprof"
	"github.com/rs/zerolog/log"
)

type Option func(config *config) error

type config struct {
	middlewares MiddlewareRegistry
}

func init() {
	defaultMiddlewareRegistry = NewMiddlewareRegistry()
	defaultMiddlewareRegistry.Register(DefaultRecovererMiddleware())
	defaultMiddlewareRegistry.Register(DefaultCorsMiddleware())
	defaultMiddlewareRegistry.Register(DefaultSecurityHeaderMiddleware())
	defaultMiddlewareRegistry.Register(DefaultRequestIDMiddleware())
	defaultMiddlewareRegistry.Register(DefaultMetricMiddleware())
}

// New creates a router with good defaults for kopexa services
func New() *chi.Mux {
	return NewWithOptions(WithMiddlewareRegistry(defaultMiddlewareRegistry))
}

func NewWithOptions(options ...Option) *chi.Mux {
	config := config{
		middlewares: NewMiddlewareRegistry(),
	}

	for _, option := range options {
		if err := option(&config); err != nil {
			log.Fatal().Err(err).Msg("error applying option")
		}
	}

	r := chi.NewRouter()
	config.middlewares.UseOnRouter(r)

	return r
}

// NewObservabilityRouter setup a HTTP server with preconfigured handlers. The support server service the prometheus
// metrics, and offers a pprof endpoint. Please don't export this SupportServer to the public internet.
func NewObservabilityRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Mount(pprof.PathPrefixPProf, pprof.Handler())
	r.Get("/metrics", metric.Handler.ServeHTTP)

	return r
}

func WithMiddlewareRegistry(registry MiddlewareRegistry) Option {
	return func(conf *config) error {
		conf.middlewares = registry
		return nil
	}
}
