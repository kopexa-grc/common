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
		field.Time("updated_at").
			Default(time.Now).
			Optional().
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
		ml, ok := m.(AuditLogger)
		if !ok {
			return nil, fmt.Errorf("%w: %T", ErrUnexpectedMutationType, m)
		}

		actor := auth.ActorFromContext(ctx)

		switch op := m.Op(); {
		case op.Is(ent.OpCreate):
			ml.SetCreatedAt(time.Now())

			if _, exists := ml.CreatedBy(); !exists {
				ml.SetCreatedBy(actor.ID)
			}
		case op.Is(ent.OpUpdateOne | ent.OpUpdate):
			ml.SetUpdatedAt(time.Now())

			if _, exists := ml.UpdatedBy(); !exists {
				ml.SetUpdatedBy(actor.ID)
			}
		}

		return next.Mutate(ctx, m)
	})
}
