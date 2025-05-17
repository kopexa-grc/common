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

// SpaceCreationContextKey is used to indicate that the current operation
// is part of a space creation process. This allows the ent privacy rules
// to bypass certain restrictions and create default entities during space setup.
//
// Example usage:
//
//	ctx = auth.WithSpaceCreation(ctx)
//	if auth.IsSpaceCreation(ctx) {
//		// Create default entities
//	}
type SpaceCreationContextKey struct{}

// WithSpaceCreation marks the context as being part of a space creation process.
// This allows the ent privacy rules to create default entities during space setup.
//
// Parameters:
//   - ctx: The context to modify
//
// Returns:
//   - context.Context: A new context with the space creation flag set
func WithSpaceCreation(ctx context.Context) context.Context {
	return context.WithValue(ctx, SpaceCreationContextKey{}, true)
}

// IsSpaceCreation checks if the current context is part of a space creation process.
//
// Parameters:
//   - ctx: The context to check
//
// Returns:
//   - bool: true if the context is part of a space creation process, false otherwise
func IsSpaceCreation(ctx context.Context) bool {
	value, ok := ctx.Value(SpaceCreationContextKey{}).(bool)
	return ok && value
}
