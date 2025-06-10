// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"errors"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

// ErrInvalidReference is returned when a reference is invalid
var ErrInvalidReference = errors.New("invalid reference")

// Reference represents a reference to a resource, which can be identified either by an ID or a KRN (Kopexa Resource Name).
// It is commonly used for referencing resources across the system.
type Reference struct {
	// ID is the unique identifier of the resource
	ID string `json:"id,omitempty"`
	// KRN is the Kopexa Resource Name of the resource
	KRN string `json:"krn,omitempty"`
}

// Validate checks if the reference is valid.
// A reference is valid if either ID or KRN is set, but not both.
//
// Returns:
//   - error: ErrInvalidReference if the reference is invalid
func (r *Reference) Validate() error {
	if r.ID == "" && r.KRN == "" {
		return fmt.Errorf("%w: either ID or KRN must be set", ErrInvalidReference)
	}

	if r.ID != "" && r.KRN != "" {
		return fmt.Errorf("%w: cannot set both ID and KRN", ErrInvalidReference)
	}

	return nil
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for Reference.
// It allows Reference to be used as a GraphQL scalar type.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (r *Reference) UnmarshalGQL(v interface{}) error {
	if err := unmarshalGQLJSON(v, r); err != nil {
		return fmt.Errorf("failed to unmarshal reference: %w", err)
	}

	return r.Validate()
}

// MarshalGQL implements the graphql.Marshaler interface for Reference.
// It allows Reference to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the Reference to
func (r Reference) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, r); err != nil {
		log.Error().Err(err).Msg("failed to marshal reference to GraphQL")
	}
}
