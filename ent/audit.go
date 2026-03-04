// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package ent

import "context"

// AuditSkipKey is a context key used to indicate that audit hook operations
// should be skipped for the current operation. This is useful when you need to
// set audit fields (created_at, created_by, updated_at, updated_by) manually,
// such as during data migrations.
type AuditSkipKey struct{}

// SkipAudit creates a new context that indicates audit hook operations
// should be skipped. This is useful when you need to manually set audit fields
// like created_at, created_by, updated_at, updated_by during migrations or
// other special operations.
//
// Example:
//
//	ctx = ent.SkipAudit(ctx)
//	client.User.Create().
//	    SetCreatedAt(customTime).
//	    SetCreatedBy(customUserID).
//	    Save(ctx) // Audit hook will not overwrite these values
func SkipAudit(parent context.Context) context.Context {
	return context.WithValue(parent, AuditSkipKey{}, true)
}

// CheckSkipAudit checks if the current context indicates that audit hook
// operations should be skipped. Returns true if audit should be skipped,
// false otherwise.
//
// Example:
//
//	if ent.CheckSkipAudit(ctx) {
//	    // Skip setting audit fields
//	} else {
//	    // Set audit fields normally
//	}
func CheckSkipAudit(ctx context.Context) bool {
	return ctx.Value(AuditSkipKey{}) != nil
}
