// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/rs/zerolog/log"
)

var (
	// ErrNilWriter is returned when a nil writer is provided
	ErrNilWriter = errors.New("writer cannot be nil")
	// ErrNilValue is returned when a nil value is provided
	ErrNilValue = errors.New("value cannot be nil")
)

// marshalGQLJSON marshals the given type into JSON and writes it to the given writer.
// It handles error cases and provides proper logging.
//
// Parameters:
//   - w: The writer to write the JSON to
//   - a: The value to marshal
//
// Returns:
//   - error: If marshaling or writing fails
func marshalGQLJSON[T any](w io.Writer, a T) error {
	if w == nil {
		return ErrNilWriter
	}

	byteData, err := json.Marshal(a)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling json object")
		return err
	}

	_, err = w.Write(byteData)
	if err != nil {
		log.Error().Err(err).Msg("error writing json object")
		return err
	}

	return nil
}

// unmarshalGQLJSON unmarshals a JSON object into the given type.
// It handles error cases and provides proper validation.
//
// Parameters:
//   - v: The value to unmarshal
//   - a: The target type to unmarshal into
//
// Returns:
//   - error: If unmarshaling fails or validation fails
func unmarshalGQLJSON[T any](v any, a T) error {
	if v == nil {
		return ErrNilValue
	}

	byteData, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteData, &a)
	if err != nil {
		return err
	}

	return nil
}
