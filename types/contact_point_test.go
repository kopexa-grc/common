// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactPoint_Validate(t *testing.T) {
	tests := []struct {
		name    string
		contact ContactPoint
		wantErr bool
	}{
		{
			name: "valid email contact",
			contact: ContactPoint{
				Method: ContactMethodEmail,
				Name:   "John Doe",
				Role:   "Support",
				Email:  "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "valid phone contact",
			contact: ContactPoint{
				Method: ContactMethodPhone,
				Name:   "Jane Smith",
				Role:   "Sales",
				Phone:  "+1234567890",
			},
			wantErr: false,
		},
		{
			name: "valid web form contact",
			contact: ContactPoint{
				Method: ContactMethodWebForm,
				Name:   "Support Team",
				URL:    "https://support.example.com",
			},
			wantErr: false,
		},
		{
			name: "missing method",
			contact: ContactPoint{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "email method without email",
			contact: ContactPoint{
				Method: ContactMethodEmail,
				Name:   "John Doe",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			contact: ContactPoint{
				Method: ContactMethodEmail,
				Name:   "John Doe",
				Email:  "invalid-email",
			},
			wantErr: true,
		},
		{
			name: "phone method without phone",
			contact: ContactPoint{
				Method: ContactMethodPhone,
				Name:   "John Doe",
			},
			wantErr: true,
		},
		{
			name: "web form method without url",
			contact: ContactPoint{
				Method: ContactMethodWebForm,
				Name:   "John Doe",
			},
			wantErr: true,
		},
		{
			name: "invalid web form url",
			contact: ContactPoint{
				Method: ContactMethodWebForm,
				Name:   "John Doe",
				URL:    "invalid-url",
			},
			wantErr: true,
		},
		{
			name: "invalid method",
			contact: ContactPoint{
				Method: "invalid",
				Name:   "John Doe",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contact.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidContactPoint)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestContactPoint_String(t *testing.T) {
	tests := []struct {
		name    string
		contact ContactPoint
		want    string
	}{
		{
			name:    "empty contact point",
			contact: ContactPoint{},
			want:    "<empty contact point>",
		},
		{
			name: "email contact with all fields",
			contact: ContactPoint{
				Method:       ContactMethodEmail,
				Name:         "John Doe",
				Role:         "Support",
				Email:        "john@example.com",
				Availability: "9-5",
			},
			want: "John Doe (Support) <john@example.com> available: 9-5",
		},
		{
			name: "phone contact with name only",
			contact: ContactPoint{
				Method: ContactMethodPhone,
				Name:   "Jane Smith",
				Phone:  "+1234567890",
			},
			want: "Jane Smith [+1234567890]",
		},
		{
			name: "web form contact with role only",
			contact: ContactPoint{
				Method: ContactMethodWebForm,
				Role:   "Support",
				URL:    "https://support.example.com",
			},
			want: "(Support) {https://support.example.com}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.contact.String())
		})
	}
}

func TestContactPoint_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   ContactPoint
		want    string
		wantErr bool
	}{
		{
			name: "valid email contact",
			input: ContactPoint{
				Method: ContactMethodEmail,
				Name:   "John Doe",
				Role:   "Support",
				Email:  "john@example.com",
			},
			want:    `{"method":"EMAIL","name":"John Doe","role":"Support","email":"john@example.com"}`,
			wantErr: false,
		},
		{
			name: "valid phone contact",
			input: ContactPoint{
				Method: ContactMethodPhone,
				Name:   "Jane Smith",
				Phone:  "+1234567890",
			},
			want:    `{"method":"PHONE","name":"Jane Smith","phone":"+1234567890"}`,
			wantErr: false,
		},
		{
			name: "valid web form contact",
			input: ContactPoint{
				Method: ContactMethodWebForm,
				URL:    "https://support.example.com",
			},
			want:    `{"method":"WEB_FORM","url":"https://support.example.com"}`,
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

			var unmarshaled ContactPoint
			err = json.Unmarshal(got, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, unmarshaled)
		})
	}
}

func TestContactPoint_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    ContactPoint
		wantErr bool
	}{
		{
			name: "valid email contact",
			input: map[string]interface{}{
				"method": "EMAIL",
				"name":   "John Doe",
				"role":   "Support",
				"email":  "john@example.com",
			},
			want: ContactPoint{
				Method: ContactMethodEmail,
				Name:   "John Doe",
				Role:   "Support",
				Email:  "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid contact - missing method",
			input: map[string]interface{}{
				"name":  "John Doe",
				"email": "john@example.com",
			},
			want:    ContactPoint{},
			wantErr: true,
		},
		{
			name: "invalid contact - invalid method",
			input: map[string]interface{}{
				"method": "invalid",
				"name":   "John Doe",
			},
			want:    ContactPoint{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ContactPoint
			err := got.UnmarshalGQL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidContactPoint)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
