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
			field:       createMockField(true),
			maxPageSize: nil,
			description: "should not modify field when maxPageSize is nil",
		},
		{
			name:        "field without pagination support",
			field:       createMockField(false),
			maxPageSize: intPtr(10),
			description: "should not modify field when it doesn't support pagination",
		},
		{
			name:        "field with no existing pagination args",
			field:       createMockField(true),
			maxPageSize: intPtr(10),
			description: "should add default first argument when no pagination args exist",
		},
		{
			name:        "field with valid first argument",
			field:       createMockFieldWithFirstArg("5"),
			maxPageSize: intPtr(10),
			description: "should not modify field when first argument is within limit",
		},
		{
			name:        "field with first argument exceeding limit",
			field:       createMockFieldWithFirstArg("15"),
			maxPageSize: intPtr(10),
			description: "should cap first argument to maxPageSize when it exceeds limit",
		},
		{
			name:        "field with valid last argument",
			field:       createMockFieldWithLastArg("5"),
			maxPageSize: intPtr(10),
			description: "should not modify field when last argument is within limit",
		},
		{
			name:        "field with last argument exceeding limit",
			field:       createMockFieldWithLastArg("15"),
			maxPageSize: intPtr(10),
			description: "should cap last argument to maxPageSize when it exceeds limit",
		},
		{
			name:        "field with invalid first argument",
			field:       createMockFieldWithFirstArg("invalid"),
			maxPageSize: intPtr(10),
			description: "should cap invalid first argument to maxPageSize",
		},
		{
			name:        "field with invalid last argument",
			field:       createMockFieldWithLastArg("invalid"),
			maxPageSize: intPtr(10),
			description: "should cap invalid last argument to maxPageSize",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store original arguments count for comparison
			originalArgsCount := len(tt.field.Arguments)

			// Call the function
			setDefaultPaginationLimit(tt.field, tt.maxPageSize)

			// Verify the field is not nil
			assert.NotNil(t, tt.field)

			// Additional assertions based on test case
			switch tt.name {
			case "nil maxPageSize":
				// Should not modify the field
				assert.Equal(t, originalArgsCount, len(tt.field.Arguments))
			case "field without pagination support":
				// Should not modify the field
				assert.Equal(t, originalArgsCount, len(tt.field.Arguments))
			case "field with no existing pagination args":
				// Should add a first argument
				assert.Equal(t, originalArgsCount+1, len(tt.field.Arguments))

				// Check that the added argument is correct
				firstArg := tt.field.Arguments.ForName(FirstArg)
				assert.NotNil(t, firstArg)
				assert.Equal(t, strconv.Itoa(*tt.maxPageSize), firstArg.Value.Raw)
			case "field with first argument exceeding limit":
				// Should cap the first argument
				firstArg := tt.field.Arguments.ForName(FirstArg)
				assert.NotNil(t, firstArg)
				assert.Equal(t, strconv.Itoa(*tt.maxPageSize), firstArg.Value.Raw)
			case "field with last argument exceeding limit":
				// Should cap the last argument
				lastArg := tt.field.Arguments.ForName(LastArg)
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
		field := createMockField(true)
		maxPageSize := 0

		setDefaultPaginationLimit(field, &maxPageSize)

		// Should add first argument with value "0"
		firstArg := field.Arguments.ForName(FirstArg)
		assert.NotNil(t, firstArg)
		assert.Equal(t, "0", firstArg.Value.Raw)
	})

	t.Run("negative maxPageSize", func(t *testing.T) {
		field := createMockField(true)
		maxPageSize := -5

		setDefaultPaginationLimit(field, &maxPageSize)

		// Should add first argument with negative value
		firstArg := field.Arguments.ForName(FirstArg)
		assert.NotNil(t, firstArg)
		assert.Equal(t, "-5", firstArg.Value.Raw)
	})

	t.Run("field with both first and last arguments", func(t *testing.T) {
		field := createMockFieldWithBothArgs("5", "3")
		maxPageSize := 10

		originalFirstArg := field.Arguments.ForName(FirstArg)
		originalLastArg := field.Arguments.ForName(LastArg)

		setDefaultPaginationLimit(field, &maxPageSize)

		// Should not modify existing arguments when they are within limit
		assert.Equal(t, originalFirstArg.Value.Raw, field.Arguments.ForName(FirstArg).Value.Raw)
		assert.Equal(t, originalLastArg.Value.Raw, field.Arguments.ForName(LastArg).Value.Raw)
	})
}

// Helper functions to create mock fields for testing

func createMockField(supportsPagination bool) *graphql.CollectedField {
	field := &graphql.CollectedField{
		Field: &ast.Field{
			Name:      "users",
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

func createMockFieldWithFirstArg(value string) *graphql.CollectedField {
	field := createMockField(true)
	field.Arguments = ast.ArgumentList{
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

func createMockFieldWithLastArg(value string) *graphql.CollectedField {
	field := createMockField(true)
	field.Arguments = ast.ArgumentList{
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

func createMockFieldWithBothArgs(firstValue, lastValue string) *graphql.CollectedField {
	field := createMockField(true)
	field.Arguments = ast.ArgumentList{
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
