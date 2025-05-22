// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

// defaultSubject defines the default type for subjects in FGA.
// If no subject type is specified, this value will be used.
const (
	defaultSubject = "user"
)

// Common error messages and constants used throughout the FGA package.
const (
	// ErrDuplicateKey is returned when attempting to write a tuple that already exists in the FGA store.
	// This is a common error that occurs when trying to create duplicate permissions.
	ErrDuplicateKey = "write a tuple which already exists"

	// Operation types for FGA tuple operations.
	// These constants define the basic operations that can be performed on FGA tuples.
	OpWrite  = "write"  // Used for creating or updating tuples
	OpDelete = "delete" // Used for removing tuples
)

// Common Permissions define standard relations that can be used across the application.
// These relations represent basic CRUD operations and are commonly used in access control.
const (
	// CanView represents the permission to view/read a resource.
	// Users with this relation can access and view the associated object.
	CanView Relation = "can_view"

	// CanEdit represents the permission to modify a resource.
	// Users with this relation can update and modify the associated object.
	CanEdit Relation = "can_edit"

	// CanDelete represents the permission to remove a resource.
	// Users with this relation can delete the associated object.
	CanDelete Relation = "can_delete"
)
