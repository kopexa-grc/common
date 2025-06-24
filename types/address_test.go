// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress_String(t *testing.T) {
	tests := []struct {
		name  string
		input Address
		want  string
	}{
		{
			name:  "empty address",
			input: Address{},
			want:  "<empty address>",
		},
		{
			name: "full address",
			input: Address{
				Line1:      "Musterstraße 123",
				Line2:      "Etage 4",
				City:       "Berlin",
				State:      "Berlin",
				PostalCode: "10115",
				Country:    "Deutschland",
			},
			want: "Musterstraße 123 Etage 4, 10115 Berlin, Deutschland",
		},
		{
			name: "address without line2",
			input: Address{
				Line1:      "Musterstraße 123",
				City:       "Berlin",
				State:      "Berlin",
				PostalCode: "10115",
				Country:    "Deutschland",
			},
			want: "Musterstraße 123, 10115 Berlin, Deutschland",
		},
		{
			name: "address with different state",
			input: Address{
				Line1:      "Musterstraße 123",
				City:       "Berlin",
				State:      "Brandenburg",
				PostalCode: "10115",
				Country:    "Deutschland",
			},
			want: "Musterstraße 123, 10115 Berlin, Brandenburg, Deutschland",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.String())
		})
	}
}

func TestAddress_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		input Address
		want  string
	}{
		{
			name:  "empty address",
			input: Address{},
			want:  `{"line1":"","line2":"","city":"","state":"","postalCode":"","country":""}`,
		},
		{
			name: "valid address",
			input: Address{
				Line1:      "Musterstraße 123",
				Line2:      "Etage 4",
				City:       "Berlin",
				State:      "Berlin",
				PostalCode: "10115",
				Country:    "Deutschland",
			},
			want: `{"line1":"Musterstraße 123","line2":"Etage 4","city":"Berlin","state":"Berlin","postalCode":"10115","country":"Deutschland"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			tt.input.MarshalGQL(&buf)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestAddress_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    Address
		wantErr bool
	}{
		{
			name: "valid address",
			input: map[string]interface{}{
				"line1":      "Musterstraße 123",
				"line2":      "Etage 4",
				"city":       "Berlin",
				"state":      "Berlin",
				"postalCode": "10115",
				"country":    "Deutschland",
			},
			want: Address{
				Line1:      "Musterstraße 123",
				Line2:      "Etage 4",
				City:       "Berlin",
				State:      "Berlin",
				PostalCode: "10115",
				Country:    "Deutschland",
			},
		},
		{
			name:    "invalid type",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var addr Address
			err := addr.UnmarshalGQL(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, addr)
		})
	}
}
