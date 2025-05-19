// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReference_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ref     Reference
		wantErr bool
	}{
		{
			name:    "valid with ID",
			ref:     Reference{ID: "123"},
			wantErr: false,
		},
		{
			name:    "valid with KRN",
			ref:     Reference{KRN: "krn:123"},
			wantErr: false,
		},
		{
			name:    "invalid empty",
			ref:     Reference{},
			wantErr: true,
		},
		{
			name:    "invalid both set",
			ref:     Reference{ID: "123", KRN: "krn:123"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ref.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidReference)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestReference_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   Reference
		want    string
		wantErr bool
	}{
		{
			name:    "with ID",
			input:   Reference{ID: "123"},
			want:    `{"id":"123"}`,
			wantErr: false,
		},
		{
			name:    "with KRN",
			input:   Reference{KRN: "krn:123"},
			want:    `{"krn":"krn:123"}`,
			wantErr: false,
		},
		{
			name:    "empty",
			input:   Reference{},
			want:    `{}`,
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

			var unmarshaled Reference
			err = json.Unmarshal(got, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, unmarshaled)
		})
	}
}

func TestReference_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    Reference
		wantErr bool
	}{
		{
			name:    "valid with ID",
			input:   map[string]interface{}{"id": "123"},
			want:    Reference{ID: "123"},
			wantErr: false,
		},
		{
			name:    "valid with KRN",
			input:   map[string]interface{}{"krn": "krn:123"},
			want:    Reference{KRN: "krn:123"},
			wantErr: false,
		},
		{
			name:    "invalid empty",
			input:   map[string]interface{}{},
			want:    Reference{},
			wantErr: true,
		},
		{
			name:    "invalid both set",
			input:   map[string]interface{}{"id": "123", "krn": "krn:123"},
			want:    Reference{ID: "123", KRN: "krn:123"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Reference
			err := got.UnmarshalGQL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidReference)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
