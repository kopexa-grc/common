// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package gql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	t.Run("FirstArg constant", func(t *testing.T) {
		assert.Equal(t, "first", FirstArg)
		assert.NotEmpty(t, FirstArg)
	})

	t.Run("LastArg constant", func(t *testing.T) {
		assert.Equal(t, "last", LastArg)
		assert.NotEmpty(t, LastArg)
	})

	t.Run("constants are different", func(t *testing.T) {
		assert.NotEqual(t, FirstArg, LastArg)
	})

	t.Run("constants are lowercase", func(t *testing.T) {
		assert.Equal(t, FirstArg, "first")
		assert.Equal(t, LastArg, "last")
	})
}
