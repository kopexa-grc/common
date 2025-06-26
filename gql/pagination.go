// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package gql provides utilities for GraphQL pagination handling.
//
// This package contains functions for managing GraphQL pagination arguments
// and enforcing pagination limits to prevent excessive data retrieval.
package gql

import (
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
)

// setDefaultPaginationLimit sets the default pagination limit for the given field.
//
// This function enforces pagination limits on GraphQL fields that support
// pagination arguments (first/last). If no pagination arguments are provided,
// it adds a default "first" argument with the maxPageSize value. If pagination
// arguments are provided but exceed the maxPageSize, they are capped to the
// maximum allowed value.
//
// The function modifies the field's arguments in place to ensure consistent
// pagination behavior across the application.
//
// Example:
//
//	setDefaultPaginationLimit(&field, &maxPageSize)
func setDefaultPaginationLimit(field *graphql.CollectedField, maxPageSize *int) {
	if field == nil || maxPageSize == nil || field.Definition == nil {
		return
	}

	defaultFirstValue := &ast.Value{
		Raw:  strconv.Itoa(*maxPageSize),
		Kind: ast.IntValue,
	}

	// Check if the field definition supports pagination arguments
	firstArgDef := field.Definition.Arguments.ForName(FirstArg)
	if firstArgDef == nil {
		return
	}

	// Check if pagination arguments are already set
	firstArg := field.Arguments.ForName(FirstArg)
	lastArg := field.Arguments.ForName(LastArg)

	// If no pagination arguments are provided, add default "first" argument
	if firstArg == nil && lastArg == nil {
		field.Arguments = append(field.Arguments, &ast.Argument{
			Name:  FirstArg,
			Value: defaultFirstValue,
		})

		return
	}

	// Validate and cap "first" argument if it exceeds maxPageSize
	if firstArg != nil && firstArg.Value != nil && firstArg.Value.Raw != "" {
		setValue, err := strconv.Atoi(firstArg.Value.Raw)
		if err != nil || setValue > *maxPageSize {
			firstArg.Value = defaultFirstValue
		}

		return
	}

	// Validate and cap "last" argument if it exceeds maxPageSize
	if lastArg != nil && lastArg.Value != nil && lastArg.Value.Raw != "" {
		setValue, err := strconv.Atoi(lastArg.Value.Raw)
		if err != nil || setValue > *maxPageSize {
			lastArg.Value = defaultFirstValue
		}

		return
	}
}
