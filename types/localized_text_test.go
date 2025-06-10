// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func TestLocalizedText_String(t *testing.T) {
	tests := []struct {
		name     string
		text     LocalizedText
		expected string
	}{
		{
			name: "English text",
			text: LocalizedText{
				Text:     "Hello",
				Language: "en",
			},
			expected: "Hello (en)",
		},
		{
			name: "German text",
			text: LocalizedText{
				Text:     "Hallo",
				Language: "de",
			},
			expected: "Hallo (de)",
		},
		{
			name: "Empty text",
			text: LocalizedText{
				Text:     "",
				Language: "en",
			},
			expected: " (en)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.text.String())
		})
	}
}

func TestLocalizedText_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		expected LocalizedText
		wantErr  bool
	}{
		{
			name:     "Simple string",
			yamlData: "Hello World",
			expected: LocalizedText{
				Text:     "Hello World",
				Language: "en",
			},
			wantErr: false,
		},
		{
			name: "Structured object",
			yamlData: `
text: Hallo Welt
language: de
`,
			expected: LocalizedText{
				Text:     "Hallo Welt",
				Language: "de",
			},
			wantErr: false,
		},
		{
			name: "Structured object with missing language",
			yamlData: `
text: Hello World
`,
			expected: LocalizedText{
				Text:     "Hello World",
				Language: "en",
			},
			wantErr: false,
		},
		{
			name:     "Invalid YAML",
			yamlData: "text: [invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result LocalizedText

			err := yaml.Unmarshal([]byte(tt.yamlData), &result)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLocalizedTextSlice_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		expected LocalizedTextSlice
		wantErr  bool
	}{
		{
			name: "Single string",
			yamlData: `
title: "Hello world"
`,
			expected: LocalizedTextSlice{
				{Text: "Hello world", Language: "en"},
			},
			wantErr: false,
		},
		{
			name: "Single mapping object",
			yamlData: `
title:
  text: "Hallo Welt"
  language: "de"
`,
			expected: LocalizedTextSlice{
				{Text: "Hallo Welt", Language: "de"},
			},
			wantErr: false,
		},
		{
			name: "List of localized texts",
			yamlData: `
title:
  - text: "Hallo"
    language: "de"
  - text: "Hello"
    language: "en"
`,
			expected: LocalizedTextSlice{
				{Text: "Hallo", Language: "de"},
				{Text: "Hello", Language: "en"},
			},
			wantErr: false,
		},
		{
			name: "List with missing language",
			yamlData: `
title:
  - text: "Bonjour"
`,
			expected: LocalizedTextSlice{
				{Text: "Bonjour", Language: "en"},
			},
			wantErr: false,
		},
		{
			name:     "Invalid YAML",
			yamlData: "title: [invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result struct {
				Title LocalizedTextSlice `yaml:"title"`
			}

			err := yaml.Unmarshal([]byte(tt.yamlData), &result)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result.Title)
		})
	}
}

func TestLocalizedTextSlice_MarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		slice    LocalizedTextSlice
		expected interface{}
	}{
		{
			name: "Single English text",
			slice: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
			},
			expected: "Hello",
		},
		{
			name: "Multiple texts",
			slice: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
				{Text: "Hallo", Language: "de"},
			},
			expected: []LocalizedText{
				{Text: "Hello", Language: "en"},
				{Text: "Hallo", Language: "de"},
			},
		},
		{
			name: "Single non-English text",
			slice: LocalizedTextSlice{
				{Text: "Hallo", Language: "de"},
			},
			expected: []LocalizedText{
				{Text: "Hallo", Language: "de"},
			},
		},
		{
			name:     "Empty slice",
			slice:    LocalizedTextSlice{},
			expected: []LocalizedText{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.slice.MarshalYAML()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLocalizedTextSlice_Value(t *testing.T) {
	tests := []struct {
		name     string
		slice    LocalizedTextSlice
		expected driver.Value
		wantErr  bool
	}{
		{
			name: "Valid slice",
			slice: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
				{Text: "Hallo", Language: "de"},
			},
			expected: []byte(`[{"text":"Hello","language":"en"},{"text":"Hallo","language":"de"}]`),
			wantErr:  false,
		},
		{
			name:     "Empty slice",
			slice:    LocalizedTextSlice{},
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.slice.Value()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLocalizedTextSlice_Equal(t *testing.T) {
	tests := []struct {
		name     string
		slice1   LocalizedTextSlice
		slice2   LocalizedTextSlice
		expected bool
	}{
		{
			name: "Equal slices",
			slice1: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
				{Text: "Hallo", Language: "de"},
			},
			slice2: LocalizedTextSlice{
				{Text: "Hallo", Language: "de"},
				{Text: "Hello", Language: "en"},
			},
			expected: true,
		},
		{
			name: "Different order",
			slice1: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
				{Text: "Hallo", Language: "de"},
			},
			slice2: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
				{Text: "Hallo", Language: "de"},
			},
			expected: true,
		},
		{
			name: "Different content",
			slice1: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
			},
			slice2: LocalizedTextSlice{
				{Text: "Hallo", Language: "de"},
			},
			expected: false,
		},
		{
			name: "Different length",
			slice1: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
			},
			slice2: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
				{Text: "Hallo", Language: "de"},
			},
			expected: false,
		},
		{
			name:     "Empty slices",
			slice1:   LocalizedTextSlice{},
			slice2:   LocalizedTextSlice{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.slice1.Equal(tt.slice2))
		})
	}
}

func TestLocalizedText_JSON(t *testing.T) {
	tests := []struct {
		name     string
		text     LocalizedText
		expected string
	}{
		{
			name: "Complete text",
			text: LocalizedText{
				Text:     "Hello",
				Language: "en",
			},
			expected: `{"text":"Hello","language":"en"}`,
		},
		{
			name: "Empty text",
			text: LocalizedText{
				Text:     "",
				Language: "en",
			},
			expected: `{"text":"","language":"en"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(jsonData))

			// Test unmarshaling
			var result LocalizedText
			err = json.Unmarshal(jsonData, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.text, result)
		})
	}
}

func TestLocalizedTextSlice_JSON(t *testing.T) {
	tests := []struct {
		name     string
		slice    LocalizedTextSlice
		expected string
	}{
		{
			name: "Multiple texts",
			slice: LocalizedTextSlice{
				{Text: "Hello", Language: "en"},
				{Text: "Hallo", Language: "de"},
			},
			expected: `[{"text":"Hello","language":"en"},{"text":"Hallo","language":"de"}]`,
		},
		{
			name:     "Empty slice",
			slice:    LocalizedTextSlice{},
			expected: `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			jsonData, err := json.Marshal(tt.slice)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(jsonData))

			// Test unmarshaling
			var result LocalizedTextSlice
			err = json.Unmarshal(jsonData, &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.slice, result)
		})
	}
}

func TestLocalizedTextSlice_GQL(t *testing.T) {
	slice := LocalizedTextSlice{
		{Text: "Hallo", Language: "de"},
		{Text: "Hello", Language: "en"},
	}

	t.Run("MarshalGQL", func(t *testing.T) {
		var buf bytes.Buffer

		slice.MarshalGQL(&buf)
		// The output should be a valid JSON array
		var out LocalizedTextSlice
		err := json.Unmarshal(buf.Bytes(), &out)
		assert.NoError(t, err)
		assert.Equal(t, slice, out)
	})

	t.Run("UnmarshalGQL", func(t *testing.T) {
		input := []map[string]string{
			{"text": "Hallo", "language": "de"},
			{"text": "Hello", "language": "en"},
		}

		var result LocalizedTextSlice
		err := result.UnmarshalGQL(input)
		assert.NoError(t, err)
		assert.Equal(t, slice, result)
	})
}
