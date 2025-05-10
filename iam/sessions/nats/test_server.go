// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package nats

import (
	"testing"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats-server/v2/test"
)

// startTestServer starts an embedded NATS server for testing
func startTestServer(t *testing.T) *server.Server {
	t.Helper()

	opts := test.DefaultTestOptions
	opts.Port = -1 // Use random port
	opts.JetStream = true

	s := test.RunServer(&opts)

	t.Cleanup(func() {
		s.Shutdown()
	})

	return s
}

// getTestServerURL returns the URL of the test server
func getTestServerURL(s *server.Server) string {
	return s.ClientURL()
}
