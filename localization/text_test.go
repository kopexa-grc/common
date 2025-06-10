// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package localization

import (
	"testing"

	"github.com/kopexa-grc/common/types"
	"github.com/stretchr/testify/assert"
)

func TestGetText(t *testing.T) {
	slice := types.LocalizedTextSlice{
		{Text: "Hallo", Language: "de"},
		{Text: "Hello", Language: "en"},
		{Text: "Bonjour", Language: "fr"},
	}

	t.Run("Exact match", func(t *testing.T) {
		assert.Equal(t, "Hallo", GetText(slice, "de"))
		assert.Equal(t, "Bonjour", GetText(slice, "fr"))
	})
	t.Run("Fallback to English", func(t *testing.T) {
		assert.Equal(t, "Hello", GetText(slice, "es"))
	})
	t.Run("Fallback to first", func(t *testing.T) {
		noEn := types.LocalizedTextSlice{
			{Text: "Hallo", Language: "de"},
			{Text: "Bonjour", Language: "fr"},
		}
		assert.Equal(t, "Hallo", GetText(noEn, "es"))
	})
	t.Run("Empty slice", func(t *testing.T) {
		assert.Equal(t, "", GetText(types.LocalizedTextSlice{}, "en"))
	})
	// No locale provided
	t.Run("No locale provided", func(t *testing.T) {
		assert.Equal(t, "Hello", GetText(slice))
	})
}

func TestGetTexts(t *testing.T) {
	slices := []types.LocalizedTextSlice{
		{{Text: "Hallo", Language: "de"}, {Text: "Hello", Language: "en"}},
		{{Text: "Welt", Language: "de"}, {Text: "World", Language: "en"}},
		{{Text: "Bonjour", Language: "fr"}, {Text: "Hello", Language: "en"}},
	}
	assert.Equal(t, []string{"Hallo", "Welt", "Hello"}, GetTexts(slices, "de"))
	assert.Equal(t, []string{"Hello", "World", "Hello"}, GetTexts(slices, "en"))
	assert.Equal(t, []string{"Hello", "World", "Hello"}, GetTexts(slices, "es"))
}

func TestHasLanguage(t *testing.T) {
	slice := types.LocalizedTextSlice{
		{Text: "Hallo", Language: "de"},
		{Text: "Hello", Language: "en"},
	}
	assert.True(t, HasLanguage(slice, "de"))
	assert.True(t, HasLanguage(slice, "en"))
	assert.False(t, HasLanguage(slice, "fr"))
	assert.False(t, HasLanguage(types.LocalizedTextSlice{}, "en"))
}

func TestGetLanguages(t *testing.T) {
	slice := types.LocalizedTextSlice{
		{Text: "Hallo", Language: "de"},
		{Text: "Hello", Language: "en"},
		{Text: "Hello", Language: "en"}, // Duplicate
		{Text: "Bonjour", Language: "fr"},
	}
	langs := GetLanguages(slice)
	assert.ElementsMatch(t, []string{"de", "en", "fr"}, langs)
	assert.Empty(t, GetLanguages(types.LocalizedTextSlice{}))
}
