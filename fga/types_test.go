// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga_test

import (
	"testing"

	"github.com/kopexa-grc/common/fga"
	openfga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/assert"
)

func TestKind_String(t *testing.T) {
	tests := []struct {
		name     string
		kind     fga.Kind
		expected string
	}{
		{
			name:     "lowercase kind",
			kind:     "user",
			expected: "user",
		},
		{
			name:     "uppercase kind",
			kind:     "USER",
			expected: "user",
		},
		{
			name:     "mixed case kind",
			kind:     "User",
			expected: "user",
		},
		{
			name:     "empty kind",
			kind:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.kind.String())
		})
	}
}

func TestRelation_String(t *testing.T) {
	tests := []struct {
		name     string
		relation fga.Relation
		expected string
	}{
		{
			name:     "lowercase relation",
			relation: "member",
			expected: "member",
		},
		{
			name:     "uppercase relation",
			relation: "MEMBER",
			expected: "member",
		},
		{
			name:     "mixed case relation",
			relation: "Member",
			expected: "member",
		},
		{
			name:     "empty relation",
			relation: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.relation.String())
		})
	}
}

func TestEntity_String(t *testing.T) {
	tests := []struct {
		name     string
		entity   fga.Entity
		expected string
	}{
		{
			name: "entity without relation",
			entity: fga.Entity{
				Kind:       "user",
				Identifier: "123",
			},
			expected: "user:123",
		},
		{
			name: "entity with relation",
			entity: fga.Entity{
				Kind:       "user",
				Identifier: "123",
				Relation:   "member",
			},
			expected: "user:123#member",
		},
		{
			name: "empty entity",
			entity: fga.Entity{
				Kind:       "",
				Identifier: "",
			},
			expected: ":",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.entity.String())
		})
	}
}

func TestParseEntity(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    fga.Entity
		expectError bool
	}{
		{
			name:  "valid entity without relation",
			input: "user:123",
			expected: fga.Entity{
				Kind:       "user",
				Identifier: "123",
			},
			expectError: false,
		},
		{
			name:  "valid entity with relation",
			input: "user:123#member",
			expected: fga.Entity{
				Kind:       "user",
				Identifier: "123",
				Relation:   "member",
			},
			expectError: false,
		},
		{
			name:        "invalid format - missing colon",
			input:       "user123",
			expected:    fga.Entity{},
			expectError: true,
		},
		{
			name:        "invalid format - multiple colons",
			input:       "user:123:456",
			expected:    fga.Entity{},
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    fga.Entity{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fga.ParseEntity(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseFGATupleKey(t *testing.T) {
	tests := []struct {
		name     string
		input    openfga.TupleKey
		expected *fga.TupleKey
	}{
		{
			name: "valid tuple key",
			input: openfga.TupleKey{
				User:     "user:123",
				Relation: "member",
				Object:   "organization:456",
			},
			expected: &fga.TupleKey{
				Subject: fga.Entity{
					Kind:       "user",
					Identifier: "123",
				},
				Relation: "member",
				Object: fga.Entity{
					Kind:       "organization",
					Identifier: "456",
				},
			},
		},
		{
			name: "invalid user format",
			input: openfga.TupleKey{
				User:     "invalid",
				Relation: "member",
				Object:   "organization:456",
			},
			expected: nil,
		},
		{
			name: "invalid object format",
			input: openfga.TupleKey{
				User:     "user:123",
				Relation: "member",
				Object:   "invalid",
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fga.ParseFGATupleKey(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
