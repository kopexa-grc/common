// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactMethod_Values(t *testing.T) {
	values := ContactMethod("").Values()
	assert.Len(t, values, 3)
	assert.Contains(t, values, string(ContactMethodEmail))
	assert.Contains(t, values, string(ContactMethodPhone))
	assert.Contains(t, values, string(ContactMethodWebForm))
}

func TestContactMethod_String(t *testing.T) {
	tests := []struct {
		name   string
		method ContactMethod
		want   string
	}{
		{
			name:   "email method",
			method: ContactMethodEmail,
			want:   "EMAIL",
		},
		{
			name:   "phone method",
			method: ContactMethodPhone,
			want:   "PHONE",
		},
		{
			name:   "web form method",
			method: ContactMethodWebForm,
			want:   "WEB_FORM",
		},
		{
			name:   "empty method",
			method: "",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.method.String())
		})
	}
}

func TestToContactMethod(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ContactMethod
	}{
		{
			name:  "valid email - uppercase",
			input: "EMAIL",
			want:  ContactMethodEmail,
		},
		{
			name:  "valid email - lowercase",
			input: "email",
			want:  ContactMethodEmail,
		},
		{
			name:  "valid email - mixed case",
			input: "Email",
			want:  ContactMethodEmail,
		},
		{
			name:  "valid phone - uppercase",
			input: "PHONE",
			want:  ContactMethodPhone,
		},
		{
			name:  "valid phone - lowercase",
			input: "phone",
			want:  ContactMethodPhone,
		},
		{
			name:  "valid web form - uppercase",
			input: "WEB_FORM",
			want:  ContactMethodWebForm,
		},
		{
			name:  "valid web form - lowercase",
			input: "web_form",
			want:  ContactMethodWebForm,
		},
		{
			name:  "empty string",
			input: "",
			want:  ContactMethodInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ToContactMethod(tt.input))
		})
	}
}

func TestContactMethod_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   ContactMethod
		want    string
		wantErr bool
	}{
		{
			name:    "email method",
			input:   ContactMethodEmail,
			want:    `"EMAIL"`,
			wantErr: false,
		},
		{
			name:    "phone method",
			input:   ContactMethodPhone,
			want:    `"PHONE"`,
			wantErr: false,
		},
		{
			name:    "web form method",
			input:   ContactMethodWebForm,
			want:    `"WEB_FORM"`,
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

			var unmarshaled ContactMethod
			err = json.Unmarshal(got, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, unmarshaled)
		})
	}
}

func TestContactMethod_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    ContactMethod
		wantErr bool
	}{
		{
			name:    "valid email method",
			input:   "EMAIL",
			want:    ContactMethodEmail,
			wantErr: false,
		},
		{
			name:    "valid phone method",
			input:   "PHONE",
			want:    ContactMethodPhone,
			wantErr: false,
		},
		{
			name:    "valid web form method",
			input:   "WEB_FORM",
			want:    ContactMethodWebForm,
			wantErr: false,
		},
		{
			name:    "wrong type - number",
			input:   123,
			want:    ContactMethodInvalid,
			wantErr: true,
		},
		{
			name:    "wrong type - boolean",
			input:   true,
			want:    ContactMethodInvalid,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ContactMethod

			err := got.UnmarshalGQL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
