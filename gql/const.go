// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package gql

// GraphQL pagination argument names.
//
// These constants define the standard GraphQL pagination argument names
// used for cursor-based pagination in GraphQL queries.
const (
	// FirstArg represents the "first" argument for forward pagination.
	// It specifies the maximum number of items to return from the beginning.
	FirstArg = "first"

	// LastArg represents the "last" argument for backward pagination.
	// It specifies the maximum number of items to return from the end.
	LastArg = "last"
)
