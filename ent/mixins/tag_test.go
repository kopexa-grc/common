// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package mixins

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTagMixin_Fields tests the Fields method of TagMixin.
//
// This test verifies that the TagMixin correctly defines the tags field
// and returns the expected number of fields.
func TestTagMixin_Fields(t *testing.T) {
	mixin := TagMixin{}
	fields := mixin.Fields()

	require.Len(t, fields, 1, "TagMixin should define exactly one field")
	assert.NotNil(t, fields[0], "field should not be nil")
}

// TestTagMixin_String tests the String method of TagMixin.
//
// This test verifies that the String method returns a meaningful
// representation of the mixin for debugging purposes.
func TestTagMixin_String(t *testing.T) {
	mixin := TagMixin{}
	result := mixin.String()

	assert.Equal(t, "TagMixin{field: tags}", result, "String should return expected representation")
}

// TestTagMixin_Embedding tests that TagMixin can be embedded in a schema.
//
// This test verifies that the mixin properly implements the ent.Mixin
// interface and can be used in actual Ent schemas.
func TestTagMixin_Embedding(t *testing.T) {
	// Create a test schema that embeds TagMixin
	type TestSchema struct {
		TagMixin
	}

	schema := TestSchema{}
	fields := schema.Fields()

	require.Len(t, fields, 1, "embedded TagMixin should provide one field")
	assert.NotNil(t, fields[0], "embedded field should not be nil")
}

// TestTagMixin_Instantiation tests that TagMixin can be instantiated.
//
// This test verifies that the mixin can be created and used without errors.
func TestTagMixin_Instantiation(t *testing.T) {
	mixin := TagMixin{}
	assert.NotNil(t, mixin, "TagMixin should be instantiable")

	fields := mixin.Fields()
	assert.NotNil(t, fields, "Fields should return a non-nil slice")
	assert.Len(t, fields, 1, "Should return exactly one field")
}

// TestTagMixin_Consistency tests that TagMixin behaves consistently.
//
// This test verifies that multiple calls to Fields() return the same result.
func TestTagMixin_Consistency(t *testing.T) {
	mixin := TagMixin{}

	fields1 := mixin.Fields()
	fields2 := mixin.Fields()

	assert.Len(t, fields1, len(fields2), "Multiple calls should return same number of fields")
	assert.Len(t, fields1, 1, "Should always return exactly one field")
}
