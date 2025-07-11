// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package mixins provides reusable schema mixins for the Ent ORM framework.
//
// This package contains common schema patterns that can be embedded into
// Ent schemas to provide consistent functionality across different entities.
// Mixins help reduce code duplication and ensure consistent behavior.
//
// Example usage:
//
//	type Document struct {
//		ent.Schema
//		mixins.TagMixin
//	}
package mixins

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// TagMixin provides a standardized tags field for Ent schemas.
//
// This mixin adds a tags field to any Ent schema that embeds it. The tags
// field is implemented as a string slice that allows for flexible categorization
// and labeling of entities. Tags are commonly used for filtering, searching,
// and organizing data.
//
// The mixin provides:
//   - A tags field with appropriate database constraints
//   - Default empty slice initialization
//   - Optional field behavior (can be null in database)
//   - Clear documentation for the field purpose
//
// Example:
//
//	type Document struct {
//		ent.Schema
//		mixins.TagMixin
//	}
//
//	// The Document will automatically have a tags field
//	// that can be used like: document.Tags = []string{"important", "draft"}
type TagMixin struct {
	mixin.Schema
}

// Fields defines the fields that this mixin adds to the schema.
//
// This method implements the ent.Mixin interface and returns the fields
// that should be added to any schema that embeds this mixin.
//
// The tags field is implemented as a string slice with the following characteristics:
//   - Field name: "tags"
//   - Type: []string (string slice)
//   - Default: Empty slice ([]string{})
//   - Optional: true (can be null in database)
//   - Comment: Descriptive documentation
//
// Returns:
//   - []ent.Field: A slice containing the tags field definition
func (t TagMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Strings("tags").
			Comment("Tags associated with the object for categorization, filtering, and organization purposes. Tags are flexible labels that can be used to group and search related entities.").
			Default([]string{}).
			Optional(),
	}
}

// String returns a string representation of the TagMixin.
//
// This method provides a human-readable representation of the mixin,
// useful for debugging and logging purposes.
func (t TagMixin) String() string {
	return "TagMixin{field: tags}"
}
