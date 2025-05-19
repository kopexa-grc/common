// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

var (
	// ErrInvalidResponseData is returned when the response data is invalid
	ErrInvalidResponseData = fmt.Errorf("invalid response data")
)

// ResponseMeta contains additional metadata about the response.
// It is a map of string keys to arbitrary values.
type ResponseMeta map[string]any

// ResponseData contains the actual answers to survey questions.
// It is a map of question keys to their corresponding values.
// The values can be of various types (string, number, boolean, array, or map).
type ResponseData map[string]ResponseDataValue

// ResponseDataValue represents a value that can be stored in a response.
// It can be one of the following types:
// - string: For text responses
// - float64: For numeric responses
// - bool: For yes/no responses
// - []string: For multiple choice responses
// - map[string]string: For key-value pair responses
// - any: For other complex responses
type ResponseDataValue any

// Merge combines the current ResponseData with another ResponseData.
// If a key exists in both maps, the value from the other map will overwrite the current value.
//
// Parameters:
//   - other: The ResponseData to merge with
func (rd *ResponseData) Merge(other ResponseData) {
	if *rd == nil {
		*rd = make(ResponseData)
	}

	for key, value := range other {
		(*rd)[key] = value
	}
}

// UnmarshalJSON implements custom JSON unmarshaling for ResponseData.
// It attempts to unmarshal values into specific types in the following order:
// 1. string
// 2. float64 (number)
// 3. bool
// 4. []string (array of strings)
// 5. map[string]string (map of strings)
// 6. interface{} (any other type)
//
// Parameters:
//   - data: The JSON data to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (rd *ResponseData) UnmarshalJSON(data []byte) error {
	// First try to unmarshal as a map[string]interface{}
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidResponseData, err)
	}

	// Initialize the map if it's nil
	if *rd == nil {
		*rd = make(ResponseData)
	}

	// Process each value in the map
	for key, rawValue := range rawMap {
		// Try to unmarshal as different types
		// Try string
		var strValue string
		if json.Unmarshal(rawValue, &strValue) == nil {
			(*rd)[key] = strValue
			continue
		}

		// Try number
		var numValue float64
		if json.Unmarshal(rawValue, &numValue) == nil {
			(*rd)[key] = numValue
			continue
		}

		// Try boolean
		var boolValue bool
		if json.Unmarshal(rawValue, &boolValue) == nil {
			(*rd)[key] = boolValue
			continue
		}

		// Try array of strings
		var arrValue []string
		if json.Unmarshal(rawValue, &arrValue) == nil {
			(*rd)[key] = arrValue
			continue
		}

		// Try map[string]string
		var mapValue map[string]string
		if json.Unmarshal(rawValue, &mapValue) == nil {
			(*rd)[key] = mapValue
			continue
		}

		// If none of the above worked, store as interface{}
		var anyValue interface{}
		if err := json.Unmarshal(rawValue, &anyValue); err != nil {
			return fmt.Errorf("%w: failed to unmarshal value for key %s: %v", ErrInvalidResponseData, key, err)
		}

		(*rd)[key] = anyValue
	}

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface for ResponseData.
// It allows ResponseData to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the ResponseData to
func (rd *ResponseData) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, rd); err != nil {
		log.Error().Err(err).Msg("failed to marshal response data to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for ResponseData.
// It allows ResponseData to be used as a GraphQL scalar type.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (rd *ResponseData) UnmarshalGQL(v interface{}) error {
	return unmarshalGQLJSON(v, rd)
}
