// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package testutils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPostgresContainer demonstrates how to use the PostgresContainer for integration tests.
func TestPostgresContainer(t *testing.T) {
	ctx := context.Background()

	container, err := NewPostgresContainer(ctx,
		WithDatabase("testdb"),
		WithUsername("testuser"),
		WithPassword("testpass"),
	)
	require.NoError(t, err)
	t.Cleanup(func() { container.Cleanup(t) })

	dsn, err := container.GetDSN(ctx)
	require.NoError(t, err)
	require.Contains(t, dsn, "host=")
	require.Contains(t, dsn, "port=")
	require.Contains(t, dsn, "user=testuser")
	require.Contains(t, dsn, "password=testpass")
	require.Contains(t, dsn, "dbname=testdb")
}
