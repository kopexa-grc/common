// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package archive

import "context"

// SkipKey is a context-key type that indicates the archive filter
// should be skipped for the current operation. The associated value is the
// boolean true.
type SkipKey struct{}

// IncludeArchived returns a new context marked to skip the archive filter.
// Use when archived objects should be included in queries or when performing
// special operations on archived records.
//
// Example:
//
//	ctx = ent.WithIncludeArchived(ctx)
//	client.User.Query().All(ctx) // Returns archived users as well
func WithIncludeArchived(parent context.Context) context.Context {
	return context.WithValue(parent, SkipKey{}, true)
}
