// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	ReadTimeout       = 5 * time.Second
	WriteTimeout      = 10 * time.Second
	IdleTimeout       = 120 * time.Second
	ReadHeaderTimeout = 5 * time.Second
	ShutdownTimeout   = 30 * time.Second
)

// CreateHTTPServer creates a http server, based on the recommendations from
// https://blog.cloudflare.com/exposing-go-on-the-internet/
func CreateHTTPServer(addr string, hand http.Handler) *http.Server {
	return &http.Server{
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
		ReadHeaderTimeout: ReadHeaderTimeout,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			CurvePreferences: []tls.CurveID{
				tls.CurveP256,
				tls.X25519,
			},
		},
		Addr:    addr,
		Handler: hand,
	}
}

// ShutdownGracefully safely shuts down http server by allowing pending processes to finish
func ShutdownGracefully(srv *http.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigCh
	log.Info().Msgf("shutting down server due to received signal: %s", sig)

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("error shutting down server")
	}

	cancel()
}
