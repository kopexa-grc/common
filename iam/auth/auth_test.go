// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package auth_test

import (
	"context"
	"testing"

	"github.com/kopexa-grc/common/iam/auth"
	"github.com/stretchr/testify/assert"
)

func TestContextFunctions(t *testing.T) {
	t.Run("WithOrganization and OrganizationFromContext", func(t *testing.T) {
		ctx := context.Background()
		orgID := "org123"

		// Test setting organization ID
		newCtx := auth.WithOrganization(ctx, orgID)
		assert.NotEqual(t, ctx, newCtx)

		// Test retrieving organization ID
		retrievedOrgID := auth.OrganizationFromContext(newCtx)
		assert.Equal(t, orgID, retrievedOrgID)

		// Test empty context
		emptyOrgID := auth.OrganizationFromContext(ctx)
		assert.Empty(t, emptyOrgID)
	})

	t.Run("WithSpace and SpaceFromContext", func(t *testing.T) {
		ctx := context.Background()
		spaceID := "space123"

		// Test setting space ID
		newCtx := auth.WithSpace(ctx, spaceID)
		assert.NotEqual(t, ctx, newCtx)

		// Test retrieving space ID
		retrievedSpaceID := auth.SpaceFromContext(newCtx)
		assert.Equal(t, spaceID, retrievedSpaceID)

		// Test empty context
		emptySpaceID := auth.SpaceFromContext(ctx)
		assert.Empty(t, emptySpaceID)
	})

	t.Run("WithActor with organization and space", func(t *testing.T) {
		ctx := context.Background()
		actor := &auth.Actor{
			ID:             "user123",
			Type:           auth.ActorTypeUser,
			Locale:         "de",
			OrganizationID: "org123",
			SpaceID:        "space123",
		}

		// Test setting actor with organization and space
		newCtx := auth.WithActor(ctx, actor)
		assert.NotEqual(t, ctx, newCtx)

		// Test retrieving actor
		retrievedActor := auth.ActorFromContext(newCtx)
		assert.Equal(t, actor, retrievedActor)

		// Test retrieving organization and space IDs
		retrievedOrgID := auth.OrganizationFromContext(newCtx)
		assert.Equal(t, actor.OrganizationID, retrievedOrgID)

		retrievedSpaceID := auth.SpaceFromContext(newCtx)
		assert.Equal(t, actor.SpaceID, retrievedSpaceID)
	})

	t.Run("WithActor without organization and space", func(t *testing.T) {
		ctx := context.Background()
		actor := &auth.Actor{
			ID:     "user123",
			Type:   auth.ActorTypeUser,
			Locale: "de",
		}

		// Test setting actor without organization and space
		newCtx := auth.WithActor(ctx, actor)
		assert.NotEqual(t, ctx, newCtx)

		// Test retrieving actor
		retrievedActor := auth.ActorFromContext(newCtx)
		assert.Equal(t, actor, retrievedActor)

		// Test retrieving organization and space IDs
		retrievedOrgID := auth.OrganizationFromContext(newCtx)
		assert.Empty(t, retrievedOrgID)

		retrievedSpaceID := auth.SpaceFromContext(newCtx)
		assert.Empty(t, retrievedSpaceID)
	})

	t.Run("Update organization with existing actor", func(t *testing.T) {
		ctx := context.Background()
		actor := &auth.Actor{
			ID:     "user123",
			Type:   auth.ActorTypeUser,
			Locale: "de",
		}

		// Set initial actor
		ctx = auth.WithActor(ctx, actor)

		// Update organization
		newOrgID := "org456"
		ctx = auth.WithOrganization(ctx, newOrgID)

		// Verify actor was updated
		retrievedActor := auth.ActorFromContext(ctx)
		assert.Equal(t, newOrgID, retrievedActor.OrganizationID)

		// Verify organization context
		retrievedOrgID := auth.OrganizationFromContext(ctx)
		assert.Equal(t, newOrgID, retrievedOrgID)
	})

	t.Run("Update space with existing actor", func(t *testing.T) {
		ctx := context.Background()
		actor := &auth.Actor{
			ID:     "user123",
			Type:   auth.ActorTypeUser,
			Locale: "de",
		}

		// Set initial actor
		ctx = auth.WithActor(ctx, actor)

		// Update space
		newSpaceID := "space456"
		ctx = auth.WithSpace(ctx, newSpaceID)

		// Verify actor was updated
		retrievedActor := auth.ActorFromContext(ctx)
		assert.Equal(t, newSpaceID, retrievedActor.SpaceID)

		// Verify space context
		retrievedSpaceID := auth.SpaceFromContext(ctx)
		assert.Equal(t, newSpaceID, retrievedSpaceID)
	})
}
