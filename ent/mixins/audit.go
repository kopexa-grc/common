// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package mixins

import (
	"context"
	"fmt"
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	cent "github.com/kopexa-grc/common/ent"
	"github.com/kopexa-grc/common/iam/auth"
)

// AuditMixin implements the ent.Mixin for sharing audit-log capabilities with package schemas.
// It provides fields for tracking creation and update timestamps and actors.
type AuditMixin struct {
	mixin.Schema
}

// Fields of the AuditMixin.
// It adds the following fields to the schema:
// - created_at: Immutable timestamp of creation
// - created_by: Optional immutable ID of the creator
// - updated_at: Timestamp of the last update, automatically updated
// - updated_by: Optional ID of the last updater
func (AuditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Immutable().
			Default(time.Now).
			Annotations(
				entgql.Skip(
					entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput,
				),
				entgql.OrderField("created_at"),
			),
		field.String("created_by").
			Immutable().
			Optional().
			Annotations(
				entgql.Skip(
					entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput,
				),
			),
		// updated_at is always set (even on create) to ensure cursor-based pagination
		// works correctly when ordering by this field. NULL values break cursor comparisons.
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Annotations(
				entgql.Skip(
					entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput,
				),
				entgql.OrderField("updated_at"),
			),
		field.String("updated_by").
			Optional().
			Annotations(
				entgql.Skip(
					entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput,
				),
			),
	}
}

// Hooks of the AuditMixin.
// It adds the AuditHook to automatically set audit fields during mutations.
func (AuditMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		AuditHook,
	}
}

// AuditHook sets and returns the created_at, updated_at, etc., fields.
// It automatically populates these fields based on the mutation operation and the actor from the context.
//
// To skip the audit hook (e.g., during migrations where you want to preserve original timestamps),
// use cent.SkipAudit(ctx) before calling the mutation.
func AuditHook(next ent.Mutator) ent.Mutator {
	type AuditLogger interface {
		SetCreatedAt(time.Time)
		CreatedAt() (value time.Time, exists bool)
		SetCreatedBy(string)
		CreatedBy() (id string, exists bool)
		SetUpdatedAt(time.Time)
		UpdatedAt() (value time.Time, exists bool)
		SetUpdatedBy(string)
		UpdatedBy() (id string, exists bool)
	}

	return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
		// Skip audit if context indicates so (useful for migrations)
		if cent.CheckSkipAudit(ctx) {
			return next.Mutate(ctx, m)
		}

		ml, ok := m.(AuditLogger)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrUnexpectedMutationType, m)
		}

		actor := auth.ActorFromContext(ctx)

		actorID := actor.ID
		if actorID == "" {
			actorID = "system"
		}

		switch op := m.Op(); {
		case op.Is(ent.OpCreate):
			if _, exists := ml.CreatedBy(); !exists {
				ml.SetCreatedBy(actorID)
			}

			if _, exists := ml.UpdatedBy(); !exists {
				ml.SetUpdatedBy(actorID)
			}
		case op.Is(ent.OpUpdateOne | ent.OpUpdate):
			if _, exists := ml.UpdatedBy(); !exists {
				ml.SetUpdatedBy(actorID)
			}
		}

		return next.Mutate(ctx, m)
	})
}
