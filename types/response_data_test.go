// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseData_Merge(t *testing.T) {
	tests := []struct {
		name     string
		initial  ResponseData
		other    ResponseData
		expected ResponseData
	}{
		{
			name:     "merge with empty",
			initial:  ResponseData{"key1": "value1"},
			other:    ResponseData{},
			expected: ResponseData{"key1": "value1"},
		},
		{
			name:     "merge into empty",
			initial:  ResponseData{},
			other:    ResponseData{"key1": "value1"},
			expected: ResponseData{"key1": "value1"},
		},
		{
			name:     "merge different keys",
			initial:  ResponseData{"key1": "value1"},
			other:    ResponseData{"key2": "value2"},
			expected: ResponseData{"key1": "value1", "key2": "value2"},
		},
		{
			name:     "merge overlapping keys",
			initial:  ResponseData{"key1": "value1"},
			other:    ResponseData{"key1": "new_value"},
			expected: ResponseData{"key1": "new_value"},
		},
		{
			name:     "merge nil map",
			initial:  nil,
			other:    ResponseData{"key1": "value1"},
			expected: ResponseData{"key1": "value1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initial.Merge(tt.other)
			assert.Equal(t, tt.expected, tt.initial)
		})
	}
}

func TestResponseData_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ResponseData
		wantErr bool
	}{
		{
			name:  "string value",
			input: `{"key1": "value1"}`,
			want: ResponseData{
				"key1": "value1",
			},
			wantErr: false,
		},
		{
			name:  "number value",
			input: `{"key1": 42}`,
			want: ResponseData{
				"key1": float64(42),
			},
			wantErr: false,
		},
		{
			name:  "boolean value",
			input: `{"key1": true}`,
			want: ResponseData{
				"key1": true,
			},
			wantErr: false,
		},
		{
			name:  "array of strings",
			input: `{"key1": ["value1", "value2"]}`,
			want: ResponseData{
				"key1": []string{"value1", "value2"},
			},
			wantErr: false,
		},
		{
			name:  "map of strings",
			input: `{"key1": {"subkey1": "value1"}}`,
			want: ResponseData{
				"key1": map[string]string{"subkey1": "value1"},
			},
			wantErr: false,
		},
		{
			name:  "mixed types",
			input: `{"key1": "value1", "key2": 42, "key3": true, "key4": ["value1"], "key5": {"subkey1": "value1"}}`,
			want: ResponseData{
				"key1": "value1",
				"key2": float64(42),
				"key3": true,
				"key4": []string{"value1"},
				"key5": map[string]string{"subkey1": "value1"},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"key1": value1}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ResponseData
			err := json.Unmarshal([]byte(tt.input), &got)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResponseData_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   ResponseData
		want    string
		wantErr bool
	}{
		{
			name: "string value",
			input: ResponseData{
				"key1": "value1",
			},
			want:    `{"key1":"value1"}`,
			wantErr: false,
		},
		{
			name: "number value",
			input: ResponseData{
				"key1": float64(42),
			},
			want:    `{"key1":42}`,
			wantErr: false,
		},
		{
			name: "boolean value",
			input: ResponseData{
				"key1": true,
			},
			want:    `{"key1":true}`,
			wantErr: false,
		},
		{
			name: "array of strings",
			input: ResponseData{
				"key1": []string{"value1", "value2"},
			},
			want:    `{"key1":["value1","value2"]}`,
			wantErr: false,
		},
		{
			name: "map of strings",
			input: ResponseData{
				"key1": map[string]string{"subkey1": "value1"},
			},
			want:    `{"key1":{"subkey1":"value1"}}`,
			wantErr: false,
		},
		{
			name: "mixed types",
			input: ResponseData{
				"key1": "value1",
				"key2": float64(42),
				"key3": true,
				"key4": []string{"value1"},
				"key5": map[string]string{"subkey1": "value1"},
			},
			want:    `{"key1":"value1","key2":42,"key3":true,"key4":["value1"],"key5":{"subkey1":"value1"}}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, string(got))

			var unmarshaled ResponseData
			err = json.Unmarshal(got, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, unmarshaled)
		})
	}
}
