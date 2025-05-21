// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package auth

import (
	"context"
	"fmt"
	"io"
)

// ContextKey is used to store values in the request context.
// It ensures type safety when retrieving context values.
type ContextKey string

const (
	// ActorContextKey is the key used to store the actor in the request context.
	ActorContextKey ContextKey = "actor"
	// OrganizationContextKey is the key used to store the organization ID in the request context.
	OrganizationContextKey ContextKey = "organization_id"
	// SpaceContextKey is the key used to store the space ID in the request context.
	SpaceContextKey ContextKey = "space_id"
)

// ActorType represents the type of actor performing an action.
// It can be either a user or a system process.
type ActorType string

const (
	// ActorTypeUser represents a human user.
	ActorTypeUser ActorType = "user"
	// ActorTypeSystem represents an automated system process.
	ActorTypeSystem ActorType = "system"
)

// Values returns a list of all possible ActorType values.
func (ActorType) Values() []string {
	return []string{
		string(ActorTypeUser),
		string(ActorTypeSystem),
	}
}

// String returns the string representation of the ActorType value.
func (r ActorType) String() string {
	return string(r)
}

// MarshalGQL implements the gqlgen Marshaler interface.
func (r ActorType) MarshalGQL(w io.Writer) {
	_, _ = w.Write([]byte(`"` + r.String() + `"`))
}

// UnmarshalGQL implements the gqlgen Unmarshaler interface.
func (r *ActorType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("wrong type for ActorType, got: %T", v) //nolint:err113
	}

	*r = ActorType(str)

	return nil
}

// Actor represents an entity that can perform actions in the system.
// It contains the ID, type, locale, and optional organization and space context.
type Actor struct {
	ID             string
	Type           ActorType
	Locale         string
	OrganizationID string
	SpaceID        string
}

// SystemActorID is the default ID used for system actors.
const SystemActorID = "system"

// WithActor stores the given actor in the context.
// It returns a new context with the actor value set.
func WithActor(ctx context.Context, actor *Actor) context.Context {
	ctx = context.WithValue(ctx, ActorContextKey, actor)
	if actor.OrganizationID != "" {
		ctx = context.WithValue(ctx, OrganizationContextKey, actor.OrganizationID)
	}

	if actor.SpaceID != "" {
		ctx = context.WithValue(ctx, SpaceContextKey, actor.SpaceID)
	}

	return ctx
}

// ActorFromContext retrieves the actor from the context.
// If no actor is found, it returns a default system actor.
func ActorFromContext(ctx context.Context) *Actor {
	if actor, ok := ctx.Value(ActorContextKey).(*Actor); ok {
		return actor
	}

	return &Actor{
		ID:   SystemActorID,
		Type: ActorTypeSystem,
	}
}

// WithOrganization stores the organization ID in the context and updates the actor if present.
// It returns a new context with the organization ID value set.
func WithOrganization(ctx context.Context, organizationID string) context.Context {
	ctx = context.WithValue(ctx, OrganizationContextKey, organizationID)

	// Update actor if present
	if actor, ok := ctx.Value(ActorContextKey).(*Actor); ok {
		actor.OrganizationID = organizationID
		ctx = context.WithValue(ctx, ActorContextKey, actor)
	}

	return ctx
}

// OrganizationFromContext retrieves the organization ID from the context.
// Returns an empty string if no organization ID is found.
func OrganizationFromContext(ctx context.Context) string {
	if orgID, ok := ctx.Value(OrganizationContextKey).(string); ok {
		return orgID
	}

	return ""
}

// WithSpace stores the space ID in the context and updates the actor if present.
// It returns a new context with the space ID value set.
func WithSpace(ctx context.Context, spaceID string) context.Context {
	ctx = context.WithValue(ctx, SpaceContextKey, spaceID)

	// Update actor if present
	if actor, ok := ctx.Value(ActorContextKey).(*Actor); ok {
		actor.SpaceID = spaceID
		ctx = context.WithValue(ctx, ActorContextKey, actor)
	}

	return ctx
}

// SpaceFromContext retrieves the space ID from the context.
// Returns an empty string if no space ID is found.
func SpaceFromContext(ctx context.Context) string {
	if spaceID, ok := ctx.Value(SpaceContextKey).(string); ok {
		return spaceID
	}

	return ""
}
