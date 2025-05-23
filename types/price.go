// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

// Price represents a monetary value with associated metadata.
// It is used to store and manipulate pricing information in a standardized format.
//
// This struct is commonly used for:
// - Subscription pricing
// - Product pricing
// - Service pricing
type Price struct {
	// Amount is the numerical value of the price (e.g., 100.00)
	// This represents the actual price value in the specified currency.
	Amount float64 `json:"amount"`

	// Interval is the time period for which the price applies (e.g., "monthly", "yearly")
	// This is used to specify the billing frequency or duration of the price.
	Interval string `json:"interval"`

	// Currency is the three-letter currency code (e.g., "USD", "EUR")
	// This specifies the currency in which the amount is denominated.
	Currency string `json:"currency"`
}

// String returns a human-readable string representation of the price.
// If the amount is 0, it returns "Free", otherwise it returns a formatted string
// containing the amount, currency, and interval.
//
// Example:
//   - Price{Amount: 0} -> "Free"
//   - Price{Amount: 100, Currency: "USD", Interval: "monthly"} -> "100(USD)/monthly"
//
// Returns:
//   - string: A formatted string representation of the price
func (p Price) String() string {
	if p.Amount == 0 {
		return "Free"
	}

	return fmt.Sprintf("%v(%s)/%s", p.Amount, p.Currency, p.Interval)
}

// MarshalGQL implements the graphql.Marshaler interface for GraphQL serialization.
// It marshals the Price struct into a JSON format suitable for GraphQL.
//
// Parameters:
//   - w: The io.Writer to write the marshaled data to
//
// The method logs any errors that occur during marshaling.
func (p Price) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, p); err != nil {
		log.Error().
			Err(err).
			Interface("price", p).
			Msg("failed to marshal price to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for GraphQL deserialization.
// It unmarshals JSON data from GraphQL into the Price struct.
//
// Parameters:
//   - v: The interface{} containing the data to unmarshal
//
// Returns:
//   - error: If the unmarshaling fails
func (p *Price) UnmarshalGQL(v interface{}) error {
	return unmarshalGQLJSON(v, p)
}
