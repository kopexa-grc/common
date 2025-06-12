// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

// RiskAudit represents a comprehensive risk assessment result containing score components,
// critical gaps, and UI hints for risk evaluation.
//
// Fields:
//   - ScoreComponents: A map of risk component names to their respective scores
//   - CriticalGaps: A list of identified critical gaps in the risk assessment
//   - UIHints: A list of hints for UI presentation of the risk assessment
type RiskAudit struct {
	ScoreComponents map[string]int `json:"scoreComponents"`
	CriticalGaps    []string       `json:"criticalGaps"`
	UIHints         []string       `json:"uiHints"`
}

// String returns a string representation of the RiskAudit.
// It formats the audit data in a human-readable format.
func (r RiskAudit) String() string {
	return fmt.Sprintf("RiskAudit{ScoreComponents: %v, CriticalGaps: %v, UIHints: %v}",
		r.ScoreComponents, r.CriticalGaps, r.UIHints)
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for RiskAudit.
// It converts the GraphQL input into a RiskAudit struct.
//
// Args:
//   - v: The interface{} value to unmarshal
//
// Returns:
//   - error: An error if unmarshaling fails
func (r *RiskAudit) UnmarshalGQL(v interface{}) error {
	if err := unmarshalGQLJSON(v, r); err != nil {
		return fmt.Errorf("failed to unmarshal reference: %w", err)
	}

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for RiskAudit.
// It converts the RiskAudit struct into a GraphQL-compatible format.
//
// Args:
//   - w: The io.Writer to write the marshaled data to
func (r RiskAudit) MarshalGQL(w io.Writer) {
	// Ensure empty slices and maps are initialized for correct JSON output
	rCopy := r
	if rCopy.ScoreComponents == nil {
		rCopy.ScoreComponents = map[string]int{}
	}

	if rCopy.CriticalGaps == nil {
		rCopy.CriticalGaps = []string{}
	}

	if rCopy.UIHints == nil {
		rCopy.UIHints = []string{}
	}

	if err := marshalGQLJSON(w, rCopy); err != nil {
		log.Error().Err(err).Msg("failed to marshal reference to GraphQL")
	}
}
