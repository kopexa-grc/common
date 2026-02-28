// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
)

// JSONContent is a wrapper around json.RawMessage that provides CSV marshaling support.
// It is used for fields that store JSON content (like Tiptap editor content) and need
// to be exported/imported via CSV.
type JSONContent json.RawMessage

// MarshalJSON implements json.Marshaler for JSONContent.
func (j JSONContent) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler for JSONContent.
func (j *JSONContent) UnmarshalJSON(data []byte) error {
	if j == nil {
		return nil
	}
	*j = JSONContent(data)
	return nil
}

// MarshalCSV implements the csvutil.Marshaler interface for JSONContent.
// It serializes the JSON content as a string for CSV export.
//
// Returns:
//   - []byte: The JSON content as a string
//   - error: If marshaling fails
func (j JSONContent) MarshalCSV() ([]byte, error) {
	if len(j) == 0 {
		return []byte(""), nil
	}
	// Return the raw JSON string
	return []byte(j), nil
}

// UnmarshalCSV implements the csvutil.Unmarshaler interface for JSONContent.
// It deserializes a JSON string from CSV import.
//
// Parameters:
//   - data: The CSV field data (JSON string)
//
// Returns:
//   - error: If unmarshaling fails
func (j *JSONContent) UnmarshalCSV(data []byte) error {
	if len(data) == 0 {
		*j = nil
		return nil
	}
	*j = JSONContent(data)
	return nil
}

// String returns the JSON content as a string.
func (j JSONContent) String() string {
	if j == nil {
		return ""
	}
	return string(j)
}
