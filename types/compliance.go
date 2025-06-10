// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"io"

	"github.com/rs/zerolog/log"
)

// ExampleEvidence represents evidence for compliance documentation.
// It is used to store information about documentation examples and their descriptions.
//
// Example:
//
//	evidence := ExampleEvidence{
//	    DocumentationType: "policy",
//	    Description: "Example of a security policy document",
//	}
type ExampleEvidence struct {
	// DocumentationType specifies the type of documentation (e.g., "policy", "procedure", "guideline")
	DocumentationType string `json:"documentationType,omitempty"`
	// Description provides a detailed explanation of the documentation example
	Description string `json:"description,omitempty"`
}

// ImplementationGuidance provides structured guidance for implementing compliance controls.
// It includes a reference to the source of the guidance and specific implementation steps.
//
// Example:
//
//	guidance := ImplementationGuidance{
//	    ReferenceID: "//kopexa.com/compliance/iso27001/2022/controls/A.5.1.1",
//	    Guidance: []string{
//	        "Implement access control policy",
//	        "Review access rights regularly",
//	    },
//	}
type ImplementationGuidance struct {
	// ReferenceID is the unique identifier for where the guidance was sourced from.
	// It should be a KRN (Kopexa Resource Name) if possible, following the format:
	// //kopexa.com/compliance/<standard>/<version>/controls/<id>
	// Example: //kopexa.com/compliance/iso27001/2022/controls/A.5.1.1
	//
	// If a KRN is not available, other reference formats can be used:
	// - ISO27001 control ID (e.g., "A.5.1.1")
	// - NIST reference (e.g., "NIST-800-53-AC-1")
	ReferenceID string `json:"referenceId,omitempty"`
	// Guidance are the steps to take to implement the control
	// Each string represents a specific action or requirement
	Guidance []string `json:"guidance,omitempty"`
}

// MarshalGQL implements the Marshaler interface for gqlgen.
// It converts the ExampleEvidence struct to a GraphQL-compatible format.
//
// Parameters:
//   - w: The writer to write the marshaled data to
func (e ExampleEvidence) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, e); err != nil {
		log.Error().Err(err).Msg("failed to marshal ExampleEvidence to GraphQL")
	}
}

// UnmarshalGQL implements the Unmarshaler interface for gqlgen.
// It converts GraphQL input into an ExampleEvidence struct.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (e *ExampleEvidence) UnmarshalGQL(v any) error {
	return unmarshalGQLJSON(v, e)
}

// MarshalGQL implements the Marshaler interface for gqlgen.
// It converts the ImplementationGuidance struct to a GraphQL-compatible format.
//
// Parameters:
//   - w: The writer to write the marshaled data to
func (i ImplementationGuidance) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, i); err != nil {
		log.Error().Err(err).Msg("failed to marshal ImplementationGuidance to GraphQL")
	}
}

// UnmarshalGQL implements the Unmarshaler interface for gqlgen.
// It converts GraphQL input into an ImplementationGuidance struct.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (i *ImplementationGuidance) UnmarshalGQL(v any) error {
	return unmarshalGQLJSON(v, i)
}
