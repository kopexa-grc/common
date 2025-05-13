// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"context"
)

// RevokeBuilder provides a fluent interface for revoking permissions.
// It allows chaining methods to construct a complete revoke request.
type RevokeBuilder struct {
	client   *Client
	subject  Entity
	object   Entity
	relation string
}

// Revoke starts a new revoke builder chain.
// Returns a new RevokeBuilder instance.
func (c *Client) Revoke() *RevokeBuilder {
	return &RevokeBuilder{
		client:  c,
		subject: Entity{Kind: "user", Identifier: ""},
	}
}

// User sets the user ID for the revoke operation.
// Returns the RevokeBuilder for method chaining.
func (b *RevokeBuilder) User(userID string) *RevokeBuilder {
	b.subject.Identifier = userID

	return b
}

// Relation sets the relation/capability to revoke.
// Returns the RevokeBuilder for method chaining.
func (b *RevokeBuilder) Relation(relation string) *RevokeBuilder {
	b.relation = relation
	return b
}

// From sets the object type and ID for the revoke operation.
// Returns the RevokeBuilder for method chaining.
func (b *RevokeBuilder) From(objectType, objectID string) *RevokeBuilder {
	b.object.Kind = Kind(objectType)
	b.object.Identifier = objectID

	return b
}

// Apply executes the revoke operation.
// Returns an error if the revoke operation fails.
func (b *RevokeBuilder) Apply(ctx context.Context) error {
	tuple := TupleKey{
		Subject:  b.subject,
		Object:   b.object,
		Relation: Relation(b.relation),
	}

	_, err := b.client.WriteTupleKeys(ctx, []TupleKey{}, []TupleKey{tuple})

	return err
}
