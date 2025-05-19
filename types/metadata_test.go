// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata_Set(t *testing.T) {
	tests := []struct {
		name     string
		initial  Metadata
		key      string
		value    string
		expected Metadata
	}{
		{
			name:     "set new key",
			initial:  Metadata{},
			key:      "key1",
			value:    "value1",
			expected: Metadata{"key1": "value1"},
		},
		{
			name:     "update existing key",
			initial:  Metadata{"key1": "old_value"},
			key:      "key1",
			value:    "new_value",
			expected: Metadata{"key1": "new_value"},
		},
		{
			name:     "set multiple keys",
			initial:  Metadata{"key1": "value1"},
			key:      "key2",
			value:    "value2",
			expected: Metadata{"key1": "value1", "key2": "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initial.Set(tt.key, tt.value)
			assert.Equal(t, tt.expected, tt.initial)
		})
	}
}

func TestMetadata_Get(t *testing.T) {
	tests := []struct {
		name        string
		metadata    Metadata
		key         string
		want        string
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "get existing key",
			metadata: Metadata{"key1": "value1"},
			key:      "key1",
			want:     "value1",
			wantErr:  false,
		},
		{
			name:        "get non-existent key",
			metadata:    Metadata{"key1": "value1"},
			key:         "key2",
			want:        "",
			wantErr:     true,
			expectedErr: ErrKeyNotFound,
		},
		{
			name:        "get from empty metadata",
			metadata:    Metadata{},
			key:         "key1",
			want:        "",
			wantErr:     true,
			expectedErr: ErrKeyNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.metadata.Get(tt.key)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)

				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMetadata_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   Metadata
		want    string
		wantErr bool
	}{
		{
			name:    "empty metadata",
			input:   Metadata{},
			want:    `{}`,
			wantErr: false,
		},
		{
			name:    "single key-value pair",
			input:   Metadata{"key1": "value1"},
			want:    `{"key1":"value1"}`,
			wantErr: false,
		},
		{
			name:    "multiple key-value pairs",
			input:   Metadata{"key1": "value1", "key2": "value2"},
			want:    `{"key1":"value1","key2":"value2"}`,
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

			var unmarshaled Metadata
			err = json.Unmarshal(got, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, unmarshaled)
		})
	}
}
