// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga

import (
	"context"
)

// GrantBuilder provides a fluent interface for granting permissions.
// It allows chaining methods to construct a complete grant request.
type GrantBuilder struct {
	client   *Client
	subject  Entity
	relation Relation
	object   Entity
}

// Grant starts a new grant builder chain.
// Returns a new GrantBuilder instance.
func (c *Client) Grant() *GrantBuilder {
	return &GrantBuilder{
		client:  c,
		subject: Entity{Kind: "user", Identifier: ""},
	}
}

// User sets the user ID for the grant.
// Returns the GrantBuilder for method chaining.
func (b *GrantBuilder) User(userID string) *GrantBuilder {
	b.subject.Identifier = userID
	b.subject.Kind = Kind("user")

	return b
}

// As sets the user type for the grant.
// Returns the GrantBuilder for method chaining.
func (b *GrantBuilder) As(userType string) *GrantBuilder {
	b.subject.Kind = Kind(userType)
	return b
}

// Relation sets the relation/capability to grant.
// Returns the GrantBuilder for method chaining.
func (b *GrantBuilder) Relation(relation string) *GrantBuilder {
	b.relation = Relation(relation)
	return b
}

// To sets the object type and ID for the grant.
// Returns the GrantBuilder for method chaining.
func (b *GrantBuilder) To(objectType, objectID string) *GrantBuilder {
	b.object.Kind = Kind(objectType)
	b.object.Identifier = objectID

	return b
}

// Apply executes the grant operation.
// Returns an error if the grant operation fails.
func (b *GrantBuilder) Apply(ctx context.Context) error {
	tuple := TupleKey{
		Subject:  b.subject,
		Object:   b.object,
		Relation: b.relation,
	}

	_, err := b.client.WriteTupleKeys(ctx, []TupleKey{tuple}, []TupleKey{})

	return err
}
