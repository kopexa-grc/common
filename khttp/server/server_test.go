// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateHTTPServer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := CreateHTTPServer(":0", handler)
	assert.NotNil(t, srv)
	assert.Equal(t, ":0", srv.Addr)
	assert.NotNil(t, srv.Handler)
	assert.Equal(t, 5*time.Second, srv.ReadTimeout)
	assert.Equal(t, 10*time.Second, srv.WriteTimeout)
	assert.Equal(t, 120*time.Second, srv.IdleTimeout)
	assert.Equal(t, 5*time.Second, srv.ReadHeaderTimeout)
	assert.NotNil(t, srv.TLSConfig)
	assert.Equal(t, uint16(0x0303), srv.TLSConfig.MinVersion) // TLS 1.2
}
