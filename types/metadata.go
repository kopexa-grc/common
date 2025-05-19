// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"errors"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

var (
	// ErrKeyNotFound is returned when a key is not found in the metadata
	ErrKeyNotFound = errors.New("key not found in metadata")
)

// Metadata represents a collection of key-value pairs for storing additional information.
// It is commonly used for storing metadata about resources, documents, or other entities.
type Metadata map[string]string

// Set adds or updates a key-value pair in the metadata.
// If the key already exists, its value will be overwritten.
//
// Parameters:
//   - key: The key to set
//   - value: The value to associate with the key
func (m Metadata) Set(key, value string) {
	m[key] = value
}

// Get retrieves a value from the metadata by its key.
// Returns an error if the key is not found.
//
// Parameters:
//   - key: The key to look up
//
// Returns:
//   - string: The value associated with the key
//   - error: ErrKeyNotFound if the key does not exist
func (m Metadata) Get(key string) (string, error) {
	value, ok := m[key]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}

	return value, nil
}

// MarshalGQL implements the graphql.Marshaler interface for Metadata.
// It allows Metadata to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the Metadata to
func (m Metadata) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, m); err != nil {
		log.Error().Err(err).Msg("failed to marshal metadata to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for Metadata.
// It allows Metadata to be used as a GraphQL scalar type.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (m *Metadata) UnmarshalGQL(v interface{}) error {
	return unmarshalGQLJSON(v, m)
}
