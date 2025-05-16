// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

// Common error messages and constants
const (
	// ErrDuplicateKey is the error message for duplicate key errors
	ErrDuplicateKey = "write a tuple which already exists"

	// Operation types
	OpWrite  = "write"
	OpDelete = "delete"
)

// Common Permissions
const (
	CanView   Relation = "can_view"
	CanEdit   Relation = "can_edit"
	CanDelete Relation = "can_delete"
)
