// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package mixins

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// IDMixin implements the ent.Mixin interface for UUID-based IDs with optional human-readable identifiers.
// It provides a robust ID system that combines:
// - A globally unique UUID as the primary identifier
// - An optional human-readable display ID with configurable prefix and length
// - Optional space-based uniqueness for both ID types
// - Built-in indexing for optimal query performance
type IDMixin struct {
	mixin.Schema

	// HumanIdentifierPrefix defines the prefix for human-readable IDs (e.g., "USR" for users).
	// If set, a display_id field will be automatically added to the schema.
	HumanIdentifierPrefix string

	// SingleFieldIndex determines if the display_id should have a unique index.
	// Set to true to enforce uniqueness across all entities using this mixin.
	SingleFieldIndex bool

	// DisplayIDLength specifies the length of the display ID without the prefix.
	// Defaults to 6 characters if not set.
	// Collision probabilities:
	// - 6 chars: ~0.005% for 10,000 IDs, ~0.5% for 100,000 IDs
	// - 8 chars: ~0.0001% for 1,000,000 IDs
	DisplayIDLength int

	// SpaceAware determines if the entity is space-aware.
	// If true, uniqueness constraints will be scoped to the space_id.
	SpaceAware bool
}

const (
	humanIDFieldName = "display_id"
	spaceIDFieldName = "space_id"
)

// Fields returns the schema fields for the IDMixin.
// It creates:
// 1. A UUID-based 'id' field as the primary identifier
// 2. Optionally, a human-readable 'display_id' field if HumanIdentifierPrefix is set
func (i IDMixin) Fields() []ent.Field {
	fields := []ent.Field{
		field.String("id").
			DefaultFunc(uuid.NewString).
			Immutable().
			Comment("Unique identifier for the entity"),
	}

	if i.HumanIdentifierPrefix != "" {
		displayField := field.String(humanIDFieldName).
			Comment(fmt.Sprintf("Human-readable identifier for the entity, prefix: %s", i.HumanIdentifierPrefix)).
			NotEmpty().
			Immutable()

		if i.SingleFieldIndex {
			displayField.Unique()
		}

		fields = append(fields, displayField)
	}

	return fields
}

// Indexes returns the schema indexes for the IDMixin.
// It ensures:
// - The 'id' field is globally unique
// - The 'display_id' field is unique if SingleFieldIndex is true
// - Space-scoped uniqueness for both 'id' and 'display_id' if SpaceAware is true
func (i IDMixin) Indexes() []ent.Index {
	idx := []ent.Index{
		index.Fields("id").
			Unique(), // enforce globally unique ids
	}

	if i.SpaceAware {
		// Add space-scoped unique index for id
		idx = append(idx, index.Fields(spaceIDFieldName, "id").
			Unique())

		if i.HumanIdentifierPrefix != "" && i.SingleFieldIndex {
			// Add space-scoped unique index for display_id
			idx = append(idx, index.Fields(spaceIDFieldName, humanIDFieldName).
				Unique())
		}
	}

	return idx
}

// Hooks returns the schema hooks for the IDMixin.
// It provides:
// - Automatic generation of human-readable IDs when HumanIdentifierPrefix is set
// - Validation and formatting of display IDs
func (i IDMixin) Hooks() []ent.Hook {
	if i.HumanIdentifierPrefix == "" {
		return []ent.Hook{}
	}

	return []ent.Hook{setIdentifierHook(i)}
}

// HookFunc defines the type for ID generation hooks
type HookFunc func(i IDMixin) ent.Hook

// setIdentifierHook creates a hook that automatically generates and sets
// human-readable display IDs based on the entity's UUID.
var setIdentifierHook HookFunc = func(i IDMixin) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			mut, ok := m.(mutationWithDisplayID)
			if ok {
				if id, exists := mut.ID(); exists {
					length := 6 // default length
					if i.DisplayIDLength > 0 {
						length = i.DisplayIDLength
					}

					out := generateShortCharID(id, length)
					mut.SetDisplayID(fmt.Sprintf("%s-%s", i.HumanIdentifierPrefix, out))
				}
			}

			return next.Mutate(ctx, m)
		})
	}
}

// generateShortCharID creates a fixed-length alphanumeric string from a UUID.
// It uses SHA256 hashing and Base32 encoding to ensure:
// - Consistent length output
// - Alphanumeric characters only (A-Z, 0-9)
// - Even distribution of values
// - Collision resistance based on length
func generateShortCharID(ulid string, length int) string {
	hash := sha256.Sum256([]byte(ulid))
	encoded := base32.StdEncoding.EncodeToString(hash[:])
	encoded = strings.ToUpper(strings.TrimRight(encoded, "="))

	return encoded[:length]
}

// mutationWithDisplayID defines the interface for mutations that support
// display ID generation and management.
type mutationWithDisplayID interface {
	SetDisplayID(string)
	ID() (id string, exists bool)
	Type() string
}
