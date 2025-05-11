// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package ent

import "context"

// SoftDeleteSkipKey is a context key used to indicate that soft delete operations
// should be skipped for the current operation. This is useful when you need to
// perform a hard delete or when you want to bypass the soft delete mechanism.
type SoftDeleteSkipKey struct{}

// SkipSoftDelete creates a new context that indicates soft delete operations
// should be skipped. This is useful when you need to perform a hard delete
// or when you want to bypass the soft delete mechanism.
//
// Example:
//
//	ctx = ent.SkipSoftDelete(ctx)
//	client.User.DeleteOne(user).Exec(ctx) // Will perform a hard delete
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, SoftDeleteSkipKey{}, true)
}

// CheckSkipSoftDelete checks if the current context indicates that soft delete
// operations should be skipped. Returns true if soft delete should be skipped,
// false otherwise.
//
// Example:
//
//	if ent.CheckSkipSoftDelete(ctx) {
//	    // Perform hard delete
//	} else {
//	    // Perform soft delete
//	}
func CheckSkipSoftDelete(ctx context.Context) bool {
	return ctx.Value(SoftDeleteSkipKey{}) != nil
}

// SoftDeleteKey is a context key used to indicate that the current operation
// is a soft delete operation. This is useful when you need to know if the
// current operation is a soft delete.
type SoftDeleteKey struct{}

// IsSoftDelete creates a new context that indicates the current operation
// is a soft delete operation. This is useful when you need to know if the
// current operation is a soft delete.
//
// Example:
//
//	ctx = ent.IsSoftDelete(ctx)
//	// Now you can check if this is a soft delete operation
func IsSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, SoftDeleteKey{}, true)
}

// CheckIsSoftDelete checks if the current context indicates that the current
// operation is a soft delete operation. Returns true if this is a soft delete
// operation, false otherwise.
//
// Example:
//
//	if ent.CheckIsSoftDelete(ctx) {
//	    // Handle soft delete specific logic
//	}
func CheckIsSoftDelete(ctx context.Context) bool {
	return ctx.Value(SoftDeleteKey{}) != nil
}
