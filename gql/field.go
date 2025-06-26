// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package gql provides utilities for GraphQL field handling and preloading.
//
// This package contains functions for analyzing GraphQL queries, extracting
// field information, and managing preloads for database operations. It supports
// nested field detection and pagination limit enforcement.
package gql

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
)

// CheckForRequestedField checks if the requested field is in the list of fields from the request.
//
// This function analyzes the GraphQL query context to determine if a specific
// field or its plural form is being requested. It performs case-insensitive
// matching and supports nested field paths.
//
// The function handles field names in the format "parent.parent.fieldName",
// e.g., "organization.orgSubscription.subscriptionURL". It also checks for
// plural versions of the field name by appending "s".
//
// Example:
//
//	if CheckForRequestedField(ctx, "user") {
//		// Handle user field request
//	}
func CheckForRequestedField(ctx context.Context, fieldName string) bool {
	// we don't care about the maxPageSize for this function
	fields := GetPreloads(ctx, nil)
	if fields == nil {
		return false
	}

	lowerFieldName := strings.ToLower(fieldName)
	pluralFieldName := lowerFieldName + "s"

	for _, f := range fields {
		lowerField := strings.ToLower(f)

		// Check if the field name is contained in the field path
		if strings.Contains(lowerField, lowerFieldName) {
			return true
		}

		// Check if it contains the plural version of the field name
		if strings.Contains(lowerField, pluralFieldName) {
			return true
		}
	}

	return false
}

// GetPreloads returns the preloads for the current GraphQL operation.
//
// This function extracts all field paths from the current GraphQL query context
// and returns them as a slice of strings. Each string represents a dot-separated
// path to a field, e.g., "organization.users.profile".
//
// If maxPageLimit is provided, it will enforce pagination limits on fields
// that support pagination arguments (first/last).
//
// Returns nil if the context is not a valid GraphQL operation context.
//
// Example:
//
//	preloads := GetPreloads(ctx, &maxPageSize)
//	for _, preload := range preloads {
//		// Process each preload path
//	}
func GetPreloads(ctx context.Context, maxPageLimit *int) []string {
	// skip if the context is not a graphql operation context
	if ok := graphql.HasOperationContext(ctx); !ok {
		return nil
	}

	gCtx := graphql.GetOperationContext(ctx)
	if gCtx == nil {
		return nil
	}

	return getNestedPreloads(
		gCtx,
		graphql.CollectFieldsCtx(ctx, nil),
		"",
		maxPageLimit,
	)
}

// getNestedPreloads returns the nested preloads for the current GraphQL operation.
//
// This function recursively traverses the GraphQL field selections to build
// a complete list of field paths. It handles nested objects and applies
// pagination limits where appropriate.
//
// The prefix parameter is used to build the full field path, starting with
// an empty string for top-level fields.
func getNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string, maxPageSize *int) (preloads []string) {
	for _, field := range fields {
		prefixField := getPreloadString(prefix, field.Name)

		// set limits on edges if max page size is set
		if maxPageSize != nil {
			setDefaultPaginationLimit(&field, maxPageSize)
		}

		// add the current field to the preloads
		preloads = append(preloads, prefixField)
		preloads = append(preloads, getNestedPreloads(ctx, graphql.CollectFields(ctx, field.Selections, nil), prefixField, maxPageSize)...)
	}

	return
}

// getPreloadString returns the preload string for the given prefix and name.
//
// This function constructs a dot-separated field path by combining the prefix
// and field name. If the prefix is empty, it returns just the field name.
//
// Example:
//   - getPreloadString("", "user") -> "user"
//   - getPreloadString("organization", "user") -> "organization.user"
func getPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}

	return name
}
