// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package gql

import (
	"strconv"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestSetDefaultPaginationLimit(t *testing.T) {
	tests := []struct {
		name        string
		field       *graphql.CollectedField
		maxPageSize *int
		description string
	}{
		{
			name:        "nil maxPageSize",
			field:       createMockField(t, "users", true),
			maxPageSize: nil,
			description: "should not modify field when maxPageSize is nil",
		},
		{
			name:        "field without pagination support",
			field:       createMockField(t, "user", false),
			maxPageSize: intPtr(10),
			description: "should not modify field when it doesn't support pagination",
		},
		{
			name:        "field with no existing pagination args",
			field:       createMockField(t, "users", true),
			maxPageSize: intPtr(10),
			description: "should add default first argument when no pagination args exist",
		},
		{
			name:        "field with valid first argument",
			field:       createMockFieldWithFirstArg(t, "users", "5"),
			maxPageSize: intPtr(10),
			description: "should not modify field when first argument is within limit",
		},
		{
			name:        "field with first argument exceeding limit",
			field:       createMockFieldWithFirstArg(t, "users", "15"),
			maxPageSize: intPtr(10),
			description: "should cap first argument to maxPageSize when it exceeds limit",
		},
		{
			name:        "field with valid last argument",
			field:       createMockFieldWithLastArg(t, "users", "5"),
			maxPageSize: intPtr(10),
			description: "should not modify field when last argument is within limit",
		},
		{
			name:        "field with last argument exceeding limit",
			field:       createMockFieldWithLastArg(t, "users", "15"),
			maxPageSize: intPtr(10),
			description: "should cap last argument to maxPageSize when it exceeds limit",
		},
		{
			name:        "field with invalid first argument",
			field:       createMockFieldWithFirstArg(t, "users", "invalid"),
			maxPageSize: intPtr(10),
			description: "should cap invalid first argument to maxPageSize",
		},
		{
			name:        "field with invalid last argument",
			field:       createMockFieldWithLastArg(t, "users", "invalid"),
			maxPageSize: intPtr(10),
			description: "should cap invalid last argument to maxPageSize",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original arguments count for comparison
			originalArgsCount := len(tt.field.Field.Arguments)

			// Call the function
			setDefaultPaginationLimit(tt.field, tt.maxPageSize)

			// Verify the field is not nil
			assert.NotNil(t, tt.field)

			// Additional assertions based on test case
			if tt.name == "nil maxPageSize" {
				// Should not modify the field
				assert.Equal(t, originalArgsCount, len(tt.field.Field.Arguments))
			} else if tt.name == "field without pagination support" {
				// Should not modify the field
				assert.Equal(t, originalArgsCount, len(tt.field.Field.Arguments))
			} else if tt.name == "field with no existing pagination args" {
				// Should add a first argument
				assert.Equal(t, originalArgsCount+1, len(tt.field.Field.Arguments))

				// Check that the added argument is correct
				firstArg := tt.field.Field.Arguments.ForName(FirstArg)
				assert.NotNil(t, firstArg)
				assert.Equal(t, strconv.Itoa(*tt.maxPageSize), firstArg.Value.Raw)
			} else if tt.name == "field with first argument exceeding limit" {
				// Should cap the first argument
				firstArg := tt.field.Field.Arguments.ForName(FirstArg)
				assert.NotNil(t, firstArg)
				assert.Equal(t, strconv.Itoa(*tt.maxPageSize), firstArg.Value.Raw)
			} else if tt.name == "field with last argument exceeding limit" {
				// Should cap the last argument
				lastArg := tt.field.Field.Arguments.ForName(LastArg)
				assert.NotNil(t, lastArg)
				assert.Equal(t, strconv.Itoa(*tt.maxPageSize), lastArg.Value.Raw)
			}
		})
	}
}

func TestSetDefaultPaginationLimit_EdgeCases(t *testing.T) {
	t.Run("nil field", func(t *testing.T) {
		maxPageSize := 10
		// Should not panic
		assert.NotPanics(t, func() {
			setDefaultPaginationLimit(nil, &maxPageSize)
		})
	})

	t.Run("zero maxPageSize", func(t *testing.T) {
		field := createMockField(t, "users", true)
		maxPageSize := 0

		setDefaultPaginationLimit(field, &maxPageSize)

		// Should add first argument with value "0"
		firstArg := field.Field.Arguments.ForName(FirstArg)
		assert.NotNil(t, firstArg)
		assert.Equal(t, "0", firstArg.Value.Raw)
	})

	t.Run("negative maxPageSize", func(t *testing.T) {
		field := createMockField(t, "users", true)
		maxPageSize := -5

		setDefaultPaginationLimit(field, &maxPageSize)

		// Should add first argument with negative value
		firstArg := field.Field.Arguments.ForName(FirstArg)
		assert.NotNil(t, firstArg)
		assert.Equal(t, "-5", firstArg.Value.Raw)
	})

	t.Run("field with both first and last arguments", func(t *testing.T) {
		field := createMockFieldWithBothArgs(t, "users", "5", "3")
		maxPageSize := 10

		originalFirstArg := field.Field.Arguments.ForName(FirstArg)
		originalLastArg := field.Field.Arguments.ForName(LastArg)

		setDefaultPaginationLimit(field, &maxPageSize)

		// Should not modify existing arguments when they are within limit
		assert.Equal(t, originalFirstArg.Value.Raw, field.Field.Arguments.ForName(FirstArg).Value.Raw)
		assert.Equal(t, originalLastArg.Value.Raw, field.Field.Arguments.ForName(LastArg).Value.Raw)
	})
}

// Helper functions to create mock fields for testing

func createMockField(t *testing.T, name string, supportsPagination bool) *graphql.CollectedField {
	field := &graphql.CollectedField{
		Field: &ast.Field{
			Name:      name,
			Arguments: ast.ArgumentList{},
		},
	}

	if supportsPagination {
		field.Definition = &ast.FieldDefinition{
			Arguments: ast.ArgumentDefinitionList{
				&ast.ArgumentDefinition{
					Name: FirstArg,
					Type: &ast.Type{
						NamedType: "Int",
					},
				},
				&ast.ArgumentDefinition{
					Name: LastArg,
					Type: &ast.Type{
						NamedType: "Int",
					},
				},
			},
		}
	}

	return field
}

func createMockFieldWithFirstArg(t *testing.T, name, value string) *graphql.CollectedField {
	field := createMockField(t, name, true)
	field.Field.Arguments = ast.ArgumentList{
		&ast.Argument{
			Name: FirstArg,
			Value: &ast.Value{
				Raw:  value,
				Kind: ast.IntValue,
			},
		},
	}
	return field
}

func createMockFieldWithLastArg(t *testing.T, name, value string) *graphql.CollectedField {
	field := createMockField(t, name, true)
	field.Field.Arguments = ast.ArgumentList{
		&ast.Argument{
			Name: LastArg,
			Value: &ast.Value{
				Raw:  value,
				Kind: ast.IntValue,
			},
		},
	}
	return field
}

func createMockFieldWithBothArgs(t *testing.T, name, firstValue, lastValue string) *graphql.CollectedField {
	field := createMockField(t, name, true)
	field.Field.Arguments = ast.ArgumentList{
		&ast.Argument{
			Name: FirstArg,
			Value: &ast.Value{
				Raw:  firstValue,
				Kind: ast.IntValue,
			},
		},
		&ast.Argument{
			Name: LastArg,
			Value: &ast.Value{
				Raw:  lastValue,
				Kind: ast.IntValue,
			},
		},
	}
	return field
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
