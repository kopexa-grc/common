// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package auth

import (
	"context"
)

// ContextKey is used to store values in the request context.
// It ensures type safety when retrieving context values.
type ContextKey string

const (
	// ActorContextKey is the key used to store the actor in the request context.
	ActorContextKey ContextKey = "actor"
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

// Actor represents an entity that can perform actions in the system.
// It contains the ID, type, and locale of the actor.
type Actor struct {
	ID     string
	Type   ActorType
	Locale string
}

// SystemActorID is the default ID used for system actors.
const SystemActorID = "system"

// WithActor stores the given actor in the context.
// It returns a new context with the actor value set.
func WithActor(ctx context.Context, actor *Actor) context.Context {
	return context.WithValue(ctx, ActorContextKey, actor)
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
