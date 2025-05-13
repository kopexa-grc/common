package fga

import "context"

// AccessBuilder provides a fluent interface for building access checks.
// It allows chaining methods to construct a complete access check request.
type AccessBuilder struct {
	client *Client
	ac     *AccessCheck
}

// Has starts a new access check builder chain.
// Returns a new AccessBuilder instance.
func (c *Client) Has() *AccessBuilder {
	return &AccessBuilder{
		client: c,
		ac:     &AccessCheck{},
	}
}

// User sets the subject ID for the access check.
// Returns the AccessBuilder for method chaining.
func (b *AccessBuilder) User(userID string) *AccessBuilder {
	b.ac.SubjectID = userID
	return b
}

// Capability sets the relation/capability to check for.
// Returns the AccessBuilder for method chaining.
func (b *AccessBuilder) Capability(capability string) *AccessBuilder {
	b.ac.Relation = capability
	return b
}

// In sets the object type and ID for the access check.
// Returns the AccessBuilder for method chaining.
func (b *AccessBuilder) In(objectType string, objectID string) *AccessBuilder {
	b.ac.ObjectType = objectType
	b.ac.ObjectID = objectID

	return b
}

// Check executes the access check and returns whether the access is granted.
// Returns true if access is granted, false otherwise, and any error that occurred.
func (b *AccessBuilder) Check(ctx context.Context) (bool, error) {
	return b.client.checkAccess(ctx, *b.ac)
}
