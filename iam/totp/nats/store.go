// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kopexa-grc/common/iam/totp"
	"github.com/nats-io/nats.go"
)

// Store implements the otpNATS interface using NATS KV
type Store struct {
	kv nats.KeyValue
}

// New creates a new NATS store
func New(js nats.JetStreamContext, bucket string) (*Store, error) {
	if js == nil {
		return nil, totp.ErrNilJetStream
	}

	kv, err := js.KeyValue(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get KV bucket: %w", err)
	}

	return &Store{
		kv: kv,
	}, nil
}

// Get retrieves a value from the store
func (s *Store) Get(_ context.Context, key string) ([]byte, error) {
	entry, err := s.kv.Get(key)
	if err != nil {
		if errors.Is(err, nats.ErrKeyNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get value: %w", err)
	}

	return entry.Value(), nil
}

// Set stores a value in the store
func (s *Store) Set(_ context.Context, key string, value []byte) error {
	_, err := s.kv.Put(key, value)
	if err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

// Delete removes a value from the store
func (s *Store) Delete(_ context.Context, key string) error {
	err := s.kv.Delete(key)
	if err != nil {
		return fmt.Errorf("failed to delete value: %w", err)
	}

	return nil
}

// GetHash retrieves a hash from the store
func (s *Store) GetHash(ctx context.Context, key string) (*totp.Hash, error) {
	data, err := s.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var hash totp.Hash
	if err := json.Unmarshal(data, &hash); err != nil {
		return nil, fmt.Errorf("failed to unmarshal hash: %w", err)
	}

	return &hash, nil
}

// SetHash stores a hash in the store
func (s *Store) SetHash(ctx context.Context, key string, hash *totp.Hash) error {
	data, err := json.Marshal(hash)
	if err != nil {
		return fmt.Errorf("failed to marshal hash: %w", err)
	}

	return s.Set(ctx, key, data)
}

// DeleteHash removes a hash from the store
func (s *Store) DeleteHash(ctx context.Context, key string) error {
	return s.Delete(ctx, key)
}

// IsExpired checks if a hash is expired
func (s *Store) IsExpired(hash *totp.Hash) bool {
	return time.Since(hash.CreatedAt) > totp.OTPExpiration*time.Minute
}
