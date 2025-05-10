// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package nats

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kopexa-grc/common/iam/totp"
)

func TestStore(t *testing.T) {
	// Create a test NATS server
	nc, err := nats.Connect(nats.DefaultURL)
	require.NoError(t, err)
	defer nc.Close()

	// Create a test JetStream context
	js, err := nc.JetStream()
	require.NoError(t, err)

	// Create a test bucket
	bucket := "test-bucket"
	_, err = js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: bucket,
	})
	require.NoError(t, err)

	// Create a new store
	store, err := New(js, bucket)
	require.NoError(t, err)

	// Test Get/Set
	t.Run("Get/Set", func(t *testing.T) {
		ctx := context.Background()
		key := "test-key"
		value := []byte("test-value")

		// Set value
		err := store.Set(ctx, key, value)
		require.NoError(t, err)

		// Get value
		got, err := store.Get(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, value, got)

		// Delete value
		err = store.Delete(ctx, key)
		require.NoError(t, err)

		// Get deleted value
		got, err = store.Get(ctx, key)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	// Test GetHash/SetHash
	t.Run("GetHash/SetHash", func(t *testing.T) {
		ctx := context.Background()
		key := "test-hash-key"
		hash := &totp.Hash{
			Hash:      "test-hash",
			CreatedAt: time.Now(),
		}

		// Set hash
		err := store.SetHash(ctx, key, hash)
		require.NoError(t, err)

		// Get hash
		got, err := store.GetHash(ctx, key)
		require.NoError(t, err)
		assert.Equal(t, hash.Hash, got.Hash)
		assert.WithinDuration(t, hash.CreatedAt, got.CreatedAt, time.Second)

		// Delete hash
		err = store.DeleteHash(ctx, key)
		require.NoError(t, err)

		// Get deleted hash
		got, err = store.GetHash(ctx, key)
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	// Test IsExpired
	t.Run("IsExpired", func(t *testing.T) {
		hash := &totp.Hash{
			Hash:      "test-hash",
			CreatedAt: time.Now().Add(-totp.OTPExpiration*time.Minute - time.Second),
		}
		assert.True(t, store.IsExpired(hash))

		hash.CreatedAt = time.Now()
		assert.False(t, store.IsExpired(hash))
	})

	// Test error cases
	t.Run("ErrorCases", func(t *testing.T) {
		// Test Get with non-existent key
		got, err := store.Get(context.Background(), "non-existent")
		require.NoError(t, err)
		assert.Nil(t, got)

		// Test GetHash with invalid JSON
		key := "invalid-json"
		err = store.Set(context.Background(), key, []byte("invalid json"))
		require.NoError(t, err)

		hash, err := store.GetHash(context.Background(), key)
		require.Error(t, err)
		assert.Nil(t, hash)

		// Test Delete with non-existent key
		err = store.Delete(context.Background(), "non-existent")
		require.NoError(t, err)
	})

	// Test New error cases
	t.Run("NewErrors", func(t *testing.T) {
		// Test with nil JetStream
		_, err := New(nil, "test-bucket")
		assert.Error(t, err)

		// Test with non-existent bucket
		_, err = New(js, "non-existent-bucket")
		assert.Error(t, err)
	})

	// Test Set error cases
	t.Run("SetErrors", func(t *testing.T) {
		// Test with nil value
		err := store.Set(context.Background(), "test-key", nil)
		assert.NoError(t, err)

		// Test with empty key
		err = store.Set(context.Background(), "", []byte("test-value"))
		assert.Error(t, err)
	})
}
