// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package models provides core data structures used throughout the application.
package types

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	DefaultCountry   = "Deutschland"
	MaxAddressLength = 100
	PostalCodeLength = 5
)

var (
	// ErrInvalidAddress is returned when an address is invalid
	ErrInvalidAddress = errors.New("invalid address")
	// ErrAddressTooLong is returned when an address field is too long
	ErrAddressTooLong = errors.New("address field too long")
)

// Address represents a physical address with all its components.
// It follows the German address format with street, house number, additional information,
// postal code, city, and country.
//
// Example:
//
//	Address{
//	  Line1: "Musterstraße 123",
//	  Line2: "Etage 4",
//	  City: "Berlin",
//	  State: "Berlin",
//	  PostalCode: "10115",
//	  Country: "Deutschland",
//	}
type Address struct {
	// Line1 is the street name and house number.
	// @example "Musterstraße 123"
	Line1 string `json:"line1" validate:"required,max=100"`

	// Line2 is additional address information (e.g., floor, apartment, company name).
	// @example "Etage 4"
	Line2 string `json:"line2" validate:"max=100"`

	// City is the city or municipality name.
	// @example "Berlin"
	City string `json:"city" validate:"required,max=100"`

	// State is the federal state (Bundesland).
	// @example "Berlin"
	State string `json:"state" validate:"required,max=100"`

	// PostalCode is the German postal code (PLZ).
	// @example "10115"
	PostalCode string `json:"postalCode" validate:"required,len=5"`

	// Country is the country name.
	// @example "Deutschland"
	Country string `json:"country" validate:"required,max=100"`
}

// Validate checks if the Address is valid.
// It verifies that all required fields are present and within length limits.
//
// Returns:
//   - error: If the Address is invalid
func (a Address) Validate() error {
	if a.Line1 == "" {
		return fmt.Errorf("%w: line1 is required", ErrInvalidAddress)
	}

	if a.City == "" {
		return fmt.Errorf("%w: city is required", ErrInvalidAddress)
	}

	if a.State == "" {
		return fmt.Errorf("%w: state is required", ErrInvalidAddress)
	}

	if a.PostalCode == "" {
		return fmt.Errorf("%w: postalCode is required", ErrInvalidAddress)
	}

	if len(a.PostalCode) != PostalCodeLength {
		return fmt.Errorf("%w: postal code must be %d digits", ErrInvalidAddress, PostalCodeLength)
	}

	// Check field lengths
	if len(a.Line1) > MaxAddressLength {
		return fmt.Errorf("%w: line1 exceeds maximum length", ErrAddressTooLong)
	}

	if len(a.Line2) > MaxAddressLength {
		return fmt.Errorf("%w: line2 exceeds maximum length", ErrAddressTooLong)
	}

	if len(a.City) > MaxAddressLength {
		return fmt.Errorf("%w: city exceeds maximum length", ErrAddressTooLong)
	}

	if len(a.State) > MaxAddressLength {
		return fmt.Errorf("%w: state exceeds maximum length", ErrAddressTooLong)
	}

	if len(a.Country) > MaxAddressLength {
		return fmt.Errorf("%w: country exceeds maximum length", ErrAddressTooLong)
	}

	return nil
}

// String returns a formatted string representation of the address.
// It combines address components in the German format, handling empty fields appropriately.
//
// Example:
//
//	"Musterstraße 123, Etage 4, 10115 Berlin, Deutschland"
//	"10115 Berlin, Deutschland" (if Line1 and Line2 are empty)
func (a Address) String() string {
	if a == (Address{}) {
		return "<empty address>"
	}

	var parts []string

	// Add street address if present
	if street := strings.TrimSpace(a.Line1 + " " + a.Line2); street != "" {
		parts = append(parts, street)
	}

	// Add postal code and city if present
	postalCity := strings.TrimSpace(a.PostalCode + " " + a.City)
	if postalCity != "" {
		parts = append(parts, postalCity)
	}

	// Add state if present and different from city
	if state := strings.TrimSpace(a.State); state != "" && state != a.City {
		parts = append(parts, state)
	}

	// Add country if present
	if country := strings.TrimSpace(a.Country); country != "" {
		parts = append(parts, country)
	}

	// If no parts are present, return empty string
	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, ", ")
}

// MarshalGQL implements the graphql.Marshaler interface for Address.
// It allows Address to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the Address to
func (a Address) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, a); err != nil {
		log.Error().Err(err).Msg("failed to marshal address to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for Address.
// It allows Address to be used as a GraphQL scalar type.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (a *Address) UnmarshalGQL(v interface{}) error {
	return unmarshalGQLJSON(v, a)
}
