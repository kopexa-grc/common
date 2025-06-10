// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package types provides core data structures and type definitions used throughout
// the application. This package includes fundamental types that are used across
// multiple domains and features.
package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog/log"
)

var (
	// ErrInvalidYAMLFormat is returned when the YAML format is invalid
	ErrInvalidYAMLFormat = errors.New("invalid YAML format for LocalizedTextSlice")
)

// LocalizedText represents a text in a specific language.
//
// This type is used to store multilingual content throughout the application.
// It provides serialization support for JSON, YAML, and GraphQL, making it
// suitable for use in various contexts including API responses, database storage,
// and configuration files.
//
// Example:
//
//	text := LocalizedText{
//		Text:     "Hello World",
//		Language: "en",
//	}
type LocalizedText struct {
	// Text contains the actual content in the specified language.
	Text string `json:"text"`

	// Language specifies the ISO language code (e.g., "en", "de", "fr").
	Language string `json:"language"`
}

// LocalizedTextSlice represents a collection of LocalizedText entries.
//
// This type is used to store multiple language versions of the same content.
// It provides methods for serialization and comparison, making it suitable for
// use in multilingual applications.
//
// Example:
//
//	texts := LocalizedTextSlice{
//		{Text: "Hello", Language: "en"},
//		{Text: "Hallo", Language: "de"},
//	}
type LocalizedTextSlice []LocalizedText

// String returns a string representation of the LocalizedText.
//
// The format is "text (language)", making it suitable for display and debugging.
//
// Example:
//
//	text := LocalizedText{Text: "Hello", Language: "en"}
//	str := text.String() // Returns "Hello (en)"
//
// Returns:
//   - string: The formatted string representation
func (l LocalizedText) String() string {
	return fmt.Sprintf("%s (%s)", l.Text, l.Language)
}

// MarshalGQL implements the graphql.Marshaler interface.
//
// This method allows LocalizedText to be used as a GraphQL scalar type.
// It serializes the LocalizedText into a JSON format that can be used in
// GraphQL responses.
//
// Parameters:
//   - w: The writer to write the serialized data to
func (l LocalizedText) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, l); err != nil {
		log.Error().Err(err).Msg("failed to marshal localized text to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface.
//
// This method allows LocalizedText to be used as a GraphQL scalar type.
// It deserializes GraphQL input into a LocalizedText structure.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (l *LocalizedText) UnmarshalGQL(v interface{}) error {
	return unmarshalGQLJSON(v, l)
}

func ToString(slice []LocalizedText, locale ...string) string {
	var fallback, english string

	targetLang := ""
	if len(locale) > 0 {
		targetLang = locale[0]
	}

	for i := range slice {
		lang := slice[i].Language
		text := slice[i].Text

		switch {
		case lang == targetLang:
			return text
		case lang == "en" && english == "":
			english = text
		case fallback == "":
			fallback = text
		}
	}

	if english != "" {
		return english
	}

	return fallback
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
//
// This method supports multiple YAML formats for LocalizedText:
//  1. Simple string format: "Hello" (defaults to English)
//  2. Structured object format:
//     text: "Hello"
//     language: "en"
//
// Example:
//
//	// Simple string format
//	var text1 LocalizedText
//	yaml.Unmarshal([]byte("Hello"), &text1)
//	// text1 = {Text: "Hello", Language: "en"}
//
//	// Structured format
//	var text2 LocalizedText
//	yaml.Unmarshal([]byte("text: Hello\nlanguage: en"), &text2)
//	// text2 = {Text: "Hello", Language: "en"}
//
// Parameters:
//   - data: The YAML data to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (l *LocalizedText) UnmarshalYAML(data []byte) error {
	// Try to unmarshal as simple string
	var raw string
	if err := yaml.Unmarshal(data, &raw); err == nil {
		l.Text = raw
		l.Language = "en"

		return nil
	}

	// Try to unmarshal as structured object
	var tmp struct {
		Text     string `yaml:"text"`
		Language string `yaml:"language"`
	}

	if err := yaml.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("invalid LocalizedText format: %w", err)
	}

	l.Text = tmp.Text
	if tmp.Language == "" {
		l.Language = "en"
	} else {
		l.Language = tmp.Language
	}

	return nil
}

func (l *LocalizedTextSlice) ToString(locale ...string) string {
	return ToString(*l, locale...)
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for LocalizedTextSlice.
//
// This method supports multiple YAML formats:
// 1. Array of LocalizedText objects
// 2. Single LocalizedText object
// 3. Array of strings (defaults to English)
//
// Example:
//
//	// Array format
//	var texts1 LocalizedTextSlice
//	yaml.Unmarshal([]byte(`
//	- text: Hello
//	  language: en
//	- text: Hallo
//	  language: de
//	`), &texts1)
//
//	// Single object format
//	var texts2 LocalizedTextSlice
//	yaml.Unmarshal([]byte(`
//	text: Hello
//	language: en
//	`), &texts2)
//
//	// String array format
//	var texts3 LocalizedTextSlice
//	yaml.Unmarshal([]byte(`
//	- Hello
//	- World
//	`), &texts3)
//
// Parameters:
//   - data: The YAML data to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (l *LocalizedTextSlice) UnmarshalYAML(data []byte) error {
	// Try to unmarshal as slice
	var list []LocalizedText
	if err := yaml.Unmarshal(data, &list); err == nil {
		for i := range list {
			if list[i].Language == "" {
				list[i].Language = "en"
			}
		}

		*l = list

		return nil
	}

	// Try to unmarshal as single structured LocalizedText
	var single LocalizedText
	if err := yaml.Unmarshal(data, &single); err == nil {
		if single.Language == "" {
			single.Language = "en"
		}

		*l = []LocalizedText{{
			Text:     single.Text,
			Language: single.Language,
		}}

		return nil
	}

	// Try to unmarshal as plain string
	var str string
	if err := yaml.Unmarshal(data, &str); err == nil {
		*l = []LocalizedText{{Text: str, Language: "en"}}
		return nil
	}

	return ErrInvalidYAMLFormat
}

// MarshalYAML implements the yaml.Marshaler interface for LocalizedTextSlice.
//
// This method provides a clean YAML representation of the LocalizedTextSlice.
// If the slice contains only one English text, it will be serialized as a
// simple string. Otherwise, it will be serialized as an array of
// LocalizedText objects.
//
// Example:
//
//	// Single English text
//	texts1 := LocalizedTextSlice{{Text: "Hello", Language: "en"}}
//	yaml.Marshal(texts1) // Returns "Hello"
//
//	// Multiple texts
//	texts2 := LocalizedTextSlice{
//		{Text: "Hello", Language: "en"},
//		{Text: "Hallo", Language: "de"},
//	}
//	yaml.Marshal(texts2) // Returns array format
//
// Returns:
//   - interface{}: The YAML-compatible representation
//   - error: If marshaling fails
func (l LocalizedTextSlice) MarshalYAML() (interface{}, error) {
	if len(l) == 1 && l[0].Language == "en" {
		return l[0].Text, nil
	}

	return []LocalizedText(l), nil
}

// Value implements the driver.Valuer interface.
//
// This method converts the LocalizedTextSlice into a format suitable for
// database storage. Empty slices are stored as NULL in the database.
//
// Returns:
//   - driver.Value: The database-compatible value
//   - error: If conversion fails
func (l LocalizedTextSlice) Value() (driver.Value, error) {
	if len(l) == 0 {
		return nil, nil
	}

	return json.Marshal(l)
}

// Equal compares two LocalizedTextSlice structures for equality.
//
// Two slices are considered equal if they contain the same texts in the
// same languages, regardless of the order of the elements.
//
// Example:
//
//	texts1 := LocalizedTextSlice{
//		{Text: "Hello", Language: "en"},
//		{Text: "Hallo", Language: "de"},
//	}
//	texts2 := LocalizedTextSlice{
//		{Text: "Hallo", Language: "de"},
//		{Text: "Hello", Language: "en"},
//	}
//	equal := texts1.Equal(texts2) // Returns true
//
// Parameters:
//   - other: The LocalizedTextSlice to compare with
//
// Returns:
//   - bool: True if the slices are equal, false otherwise
func (l LocalizedTextSlice) Equal(other LocalizedTextSlice) bool {
	if len(l) != len(other) {
		return false
	}

	lMap := make(map[string]string)
	otherMap := make(map[string]string)

	for i := range l {
		lMap[l[i].Language] = l[i].Text
	}

	for i := range other {
		otherMap[other[i].Language] = other[i].Text
	}

	for lang, text := range lMap {
		if otherText, exists := otherMap[lang]; !exists || otherText != text {
			return false
		}
	}

	return true
}

// MarshalGQL implements the graphql.Marshaler interface for LocalizedTextSlice.
//
// This method allows LocalizedTextSlice to be used as a GraphQL scalar type.
// It serializes the slice into a JSON array for GraphQL responses.
//
// Parameters:
//   - w: The writer to write the serialized data to
func (l LocalizedTextSlice) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, l); err != nil {
		log.Error().Err(err).Msg("failed to marshal LocalizedTextSlice to GraphQL")
	}
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for LocalizedTextSlice.
//
// This method allows LocalizedTextSlice to be used as a GraphQL scalar type.
// It deserializes GraphQL input into a LocalizedTextSlice structure.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (l *LocalizedTextSlice) UnmarshalGQL(v interface{}) error {
	return unmarshalGQLJSON(v, l)
}
