// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package nats

import (
	"context"
	"testing"
	"time"

	"github.com/kopexa-grc/common/iam/totp"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore(t *testing.T) {
	// Start embedded NATS server with JetStream enabled
	opts := &server.Options{
		Port:      -1, // Random port
		JetStream: true,
	}

	s := test.RunServer(opts)
	defer s.Shutdown()

	// Connect to NATS
	nc, err := nats.Connect(s.ClientURL())
	require.NoError(t, err)
	defer nc.Close()

	// Create JetStream context
	js, err := nc.JetStream()
	require.NoError(t, err)

	// Create KV bucket
	_, err = js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "totp",
	})
	require.NoError(t, err)

	// Create store
	store, err := New(js, "totp")
	require.NoError(t, err)

	// Test cases
	t.Run("SetAndGet", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"
		value := []byte("test-value")

		err := store.Set(ctx, key, value)
		assert.NoError(t, err)

		got, err := store.Get(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, value, got)
	})

	t.Run("GetNonExistent", func(t *testing.T) {
		ctx := context.Background()
		key := "non-existent"

		got, err := store.Get(ctx, key)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("Delete", func(t *testing.T) {
		ctx := context.Background()
		key := "delete-key"
		value := []byte("delete-value")

		err := store.Set(ctx, key, value)
		assert.NoError(t, err)

		err = store.Delete(ctx, key)
		assert.NoError(t, err)

		got, err := store.Get(ctx, key)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("SetAndGetHash", func(t *testing.T) {
		ctx := context.Background()
		key := "hash-key"
		hash := &totp.Hash{
			Hash:      "test-hash",
			CreatedAt: time.Now(),
		}

		err := store.SetHash(ctx, key, hash)
		assert.NoError(t, err)

		got, err := store.GetHash(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, hash.Hash, got.Hash)
		assert.WithinDuration(t, hash.CreatedAt, got.CreatedAt, time.Second)
	})

	t.Run("DeleteHash", func(t *testing.T) {
		ctx := context.Background()
		key := "delete-hash-key"
		hash := &totp.Hash{
			Hash:      "test-hash",
			CreatedAt: time.Now(),
		}

		err := store.SetHash(ctx, key, hash)
		assert.NoError(t, err)

		err = store.DeleteHash(ctx, key)
		assert.NoError(t, err)

		got, err := store.GetHash(ctx, key)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("IsExpired", func(t *testing.T) {
		hash := &totp.Hash{
			Hash:      "test-hash",
			CreatedAt: time.Now().Add(-totp.OTPExpiration*time.Minute - time.Second),
		}
		assert.True(t, store.IsExpired(hash))

		hash = &totp.Hash{
			Hash:      "test-hash",
			CreatedAt: time.Now(),
		}
		assert.False(t, store.IsExpired(hash))
	})
}
