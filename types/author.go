// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	// ErrInvalidAuthor is returned when an author is invalid
	ErrInvalidAuthor = errors.New("invalid author")
)

// Author represents a person who has created or modified a resource.
// It contains basic information about the author including their ID, name, and email.
type Author struct {
	// ID is the unique identifier of the author
	ID string `json:"id" yaml:"id"`
	// Name is the display name of the author
	Name string `json:"name" yaml:"name"`
	// Email is the email address of the author
	Email string `json:"email" yaml:"email"`
}

// Validate checks if the author is valid.
// An author is valid if it has a non-empty name and a valid email address.
//
// Returns:
//   - error: ErrInvalidAuthor if the author is invalid
func (a Author) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidAuthor)
	}

	if a.Email == "" {
		return fmt.Errorf("%w: email is required", ErrInvalidAuthor)
	}

	if !strings.Contains(a.Email, "@") {
		return fmt.Errorf("%w: invalid email format", ErrInvalidAuthor)
	}

	return nil
}

// String returns a string representation of the author.
// If the author is empty, it returns "<empty author>".
// Otherwise, it returns the author's name and email in the format "name <email>".
func (a Author) String() string {
	if a == (Author{}) {
		return "<empty author>"
	}

	return fmt.Sprintf("%s <%s>", a.Name, a.Email)
}

// MarshalGQL implements the graphql.Marshaler interface for Author.
// It allows Author to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the Author to
func (a Author) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, a); err != nil {
		log.Error().Err(err).Msg("failed to marshal author to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for Author.
// It allows Author to be used as a GraphQL scalar type.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (a *Author) UnmarshalGQL(v interface{}) error {
	if err := unmarshalGQLJSON(v, a); err != nil {
		return fmt.Errorf("failed to unmarshal author: %w", err)
	}

	return nil
}
