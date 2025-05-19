// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthor_Validate(t *testing.T) {
	tests := []struct {
		name    string
		author  Author
		wantErr bool
	}{
		{
			name: "valid author",
			author: Author{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			author: Author{
				Email: "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			author: Author{
				Name: "John Doe",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			author: Author{
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: true,
		},
		{
			name:    "empty author",
			author:  Author{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.author.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidAuthor)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestAuthor_String(t *testing.T) {
	tests := []struct {
		name   string
		author Author
		want   string
	}{
		{
			name: "valid author",
			author: Author{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			want: "John Doe <john@example.com>",
		},
		{
			name:   "empty author",
			author: Author{},
			want:   "<empty author>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.author.String())
		})
	}
}

func TestAuthor_JSON(t *testing.T) {
	tests := []struct {
		name    string
		input   Author
		want    string
		wantErr bool
	}{
		{
			name: "valid author",
			input: Author{
				ID:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			want:    `{"id":"123","name":"John Doe","email":"john@example.com"}`,
			wantErr: false,
		},
		{
			name: "minimal author",
			input: Author{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			want:    `{"id":"","name":"John Doe","email":"john@example.com"}`,
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

			var unmarshaled Author
			err = json.Unmarshal(got, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, unmarshaled)
		})
	}
}

func TestAuthor_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    Author
		wantErr bool
	}{
		{
			name: "valid author",
			input: map[string]interface{}{
				"id":    "123",
				"name":  "John Doe",
				"email": "john@example.com",
			},
			want: Author{
				ID:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid author - missing name",
			input: map[string]interface{}{
				"email": "john@example.com",
			},
			want:    Author{},
			wantErr: true,
		},
		{
			name: "invalid author - missing email",
			input: map[string]interface{}{
				"name": "John Doe",
			},
			want:    Author{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Author
			err := got.UnmarshalGQL(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalidAuthor)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
