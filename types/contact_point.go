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
	// ErrInvalidContactPoint is returned when a contact point is invalid
	ErrInvalidContactPoint = errors.New("invalid contact point")
)

// ContactPoint contains vendor contact information.
// It represents a single point of contact with a specific method and details.
type ContactPoint struct {
	// Method is the primary contact method (e.g., "email", "phone", "web")
	Method ContactMethod `json:"method" yaml:"method" validate:"required"`
	// Name is the name of the contact person
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Role is the role or position of the contact person
	Role string `json:"role,omitempty" yaml:"role,omitempty"`
	// Email is the email address for contact
	Email string `json:"email,omitempty" yaml:"email,omitempty" validate:"email"`
	// Phone is the phone number for contact
	Phone string `json:"phone,omitempty" yaml:"phone,omitempty"`
	// URL is the web URL for contact
	URL string `json:"url,omitempty" yaml:"url,omitempty"`
	// Availability describes when the contact is available
	Availability string `json:"availability,omitempty" yaml:"availability,omitempty"`
}

// Validate checks if the contact point is valid.
// A contact point is valid if:
// - It has a valid method
// - If method is email, it has a valid email address
// - If method is phone, it has a phone number
// - If method is web form, it has a URL
//
// Returns:
//   - error: ErrInvalidContactPoint if the contact point is invalid
func (c ContactPoint) Validate() error {
	if c.Method == "" {
		return fmt.Errorf("%w: method is required", ErrInvalidContactPoint)
	}

	switch c.Method {
	case ContactMethodEmail:
		if c.Email == "" {
			return fmt.Errorf("%w: email is required for email contact method", ErrInvalidContactPoint)
		}
		if !strings.Contains(c.Email, "@") {
			return fmt.Errorf("%w: invalid email format", ErrInvalidContactPoint)
		}
	case ContactMethodPhone:
		if c.Phone == "" {
			return fmt.Errorf("%w: phone is required for phone contact method", ErrInvalidContactPoint)
		}
	case ContactMethodWebForm:
		if c.URL == "" {
			return fmt.Errorf("%w: url is required for web form contact method", ErrInvalidContactPoint)
		}
		if !strings.HasPrefix(c.URL, "http://") && !strings.HasPrefix(c.URL, "https://") {
			return fmt.Errorf("%w: url must start with http:// or https://", ErrInvalidContactPoint)
		}
	default:
		return fmt.Errorf("%w: invalid contact method: %s", ErrInvalidContactPoint, c.Method)
	}

	return nil
}

// String returns a string representation of the contact point.
// If the contact point is empty, it returns "<empty contact point>".
// Otherwise, it returns the contact details in a readable format.
func (c ContactPoint) String() string {
	if c == (ContactPoint{}) {
		return "<empty contact point>"
	}

	var details []string
	if c.Name != "" {
		details = append(details, c.Name)
	}
	if c.Role != "" {
		details = append(details, fmt.Sprintf("(%s)", c.Role))
	}

	switch c.Method {
	case ContactMethodEmail:
		details = append(details, fmt.Sprintf("<%s>", c.Email))
	case ContactMethodPhone:
		details = append(details, fmt.Sprintf("[%s]", c.Phone))
	case ContactMethodWebForm:
		details = append(details, fmt.Sprintf("{%s}", c.URL))
	}

	if c.Availability != "" {
		details = append(details, fmt.Sprintf("available: %s", c.Availability))
	}

	return strings.Join(details, " ")
}

// MarshalGQL implements the graphql.Marshaler interface for ContactPoint.
// It allows ContactPoint to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the ContactPoint to
func (c ContactPoint) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, c); err != nil {
		log.Error().Err(err).Msg("failed to marshal contact point to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for ContactPoint.
// It allows ContactPoint to be used as a GraphQL scalar type.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (c *ContactPoint) UnmarshalGQL(v interface{}) error {
	if err := unmarshalGQLJSON(v, c); err != nil {
		return fmt.Errorf("failed to unmarshal contact point: %w", err)
	}
	return c.Validate()
}
