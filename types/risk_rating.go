// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

const (
	// MinRiskValue is the minimum allowed value for likelihood and consequence
	MinRiskValue = 1
	// MaxRiskValue is the maximum allowed value for likelihood and consequence
	MaxRiskValue = 5
	// DefaultScore is returned when the risk rating is zero
	DefaultScore = 25
)

// ErrInvalidRiskValue is returned when likelihood or consequence is outside valid range
var ErrInvalidRiskValue = fmt.Errorf("risk value must be between %d and %d", MinRiskValue, MaxRiskValue)

// ActionRisk represents the risk assessment details for an action.
// It includes both the inherent risk (before treatment) and the residual risk (after treatment).
type RiskRating struct {
	// Inherent risk assessment

	// Likelihood represents the probability of the risk occurring on a scale of 1-5.
	// Higher values indicate greater probability.
	Likelihood int `json:"likelihood" example:"3"`

	// Consequence represents the impact if the risk occurs on a scale of 1-5.
	// Higher values indicate more severe consequences.
	Consequence int `json:"consequence" example:"4"`

	// Rating is calculated as Likelihood × Consequence, representing the overall risk level.
	// This is automatically calculated and should not be set directly.
	Rating int `json:"rating" example:"12"`

	// Comment provides additional context or explanation for the risk assessment.
	Comment string `json:"comment,omitempty" example:"This is a comment"`
}

// CalculateRatings computes the risk rating based on likelihood and consequence.
// The rating is calculated as Likelihood × Consequence.
func (r *RiskRating) CalculateRatings() {
	r.Rating = r.Likelihood * r.Consequence
}

// IsComplete checks if the risk assessment data is complete.
// Returns true if both likelihood and consequence are greater than 0.
func (r *RiskRating) IsComplete() bool {
	return r.Likelihood > 0 && r.Consequence > 0
}

// IsZero checks if the risk rating is unset (both likelihood and consequence are 0).
// Returns true if both values are 0.
func (r *RiskRating) IsZero() bool {
	return r.Likelihood == 0 && r.Consequence == 0
}

// IsInvalid checks if the risk rating values are outside the valid range.
// Returns true if either likelihood or consequence is less than MinRiskValue or greater than MaxRiskValue.
func (r *RiskRating) IsInvalid() bool {
	return r.Likelihood < MinRiskValue || r.Likelihood > MaxRiskValue ||
		r.Consequence < MinRiskValue || r.Consequence > MaxRiskValue
}

// Score calculates the risk score based on likelihood and consequence.
// Returns DefaultScore if the risk rating is zero, otherwise returns Likelihood × Consequence.
func (r *RiskRating) Score() int {
	if r.IsZero() {
		return DefaultScore
	}

	return r.Likelihood * r.Consequence
}

// String returns a string representation of the risk rating.
// Format: "L{likelihood}C{consequence}"
func (r *RiskRating) String() string {
	return fmt.Sprintf("L%dC%d", r.Likelihood, r.Consequence)
}

// MarshalGQL implements the graphql.Marshaler interface for RiskRating.
// It allows RiskRating to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the RiskRating to
func (r RiskRating) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, r); err != nil {
		log.Error().Err(err).Msg("failed to marshal risk rating to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for RiskRating.
// It allows RiskRating to be used as a GraphQL scalar type.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (r *RiskRating) UnmarshalGQL(v interface{}) error {
	return unmarshalGQLJSON(v, r)
}
