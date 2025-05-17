// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package auth

import "context"

// OrganizationCreationContextKey is used to indicate that the current operation
// is part of an organization creation process. This allows the ent privacy rules
// to bypass certain restrictions and create default entities during organization setup.
//
// Example usage:
//
//	ctx = auth.WithOrganizationCreation(ctx)
//	if auth.IsOrganizationCreation(ctx) {
//		// Create default entities
//	}
type OrganizationCreationContextKey struct{}

// WithOrganizationCreation marks the context as being part of an organization creation process.
// This allows the ent privacy rules to create default entities during organization setup.
//
// Parameters:
//   - ctx: The context to modify
//
// Returns:
//   - context.Context: A new context with the organization creation flag set
func WithOrganizationCreation(ctx context.Context) context.Context {
	return context.WithValue(ctx, OrganizationCreationContextKey{}, true)
}

// IsOrganizationCreation checks if the current context is part of an organization creation process.
//
// Parameters:
//   - ctx: The context to check
//
// Returns:
//   - bool: true if the context is part of an organization creation process, false otherwise
func IsOrganizationCreation(ctx context.Context) bool {
	value, ok := ctx.Value(OrganizationCreationContextKey{}).(bool)
	return ok && value
}
