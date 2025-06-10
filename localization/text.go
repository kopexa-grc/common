// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package localization provides utilities for handling localized text content.
//
// This package implements a flexible text localization system that supports
// multiple languages and fallback mechanisms. It is designed to work with
// the types.LocalizedText and types.LocalizedTextSlice types.
package localization

import (
	"github.com/kopexa-grc/common/types"
)

// GetText retrieves the most appropriate text for a given locale from a
// LocalizedTextSlice.
//
// The function implements a three-tier fallback mechanism:
// 1. Returns the text in the requested locale if available
// 2. Falls back to English if the requested locale is not available
// 3. Uses the first available text as a last resort
//
// Example:
//
//	texts := types.LocalizedTextSlice{
//		{Text: "Hallo", Language: "de"},
//		{Text: "Hello", Language: "en"},
//		{Text: "Bonjour", Language: "fr"},
//	}
//	text := GetText(texts, "de") // Returns "Hallo"
//	text = GetText(texts, "es")  // Returns "Hello" (English fallback)
//
// Parameters:
//   - slice: The LocalizedTextSlice containing the localized texts
//   - locale: Optional locale code (e.g., "en", "de", "fr"). If not provided,
//     the function will return the first available text.
//
// Returns:
//   - string: The most appropriate text for the given locale
func GetText(slice types.LocalizedTextSlice, locale ...string) string {
	var fallback, english string

	targetLang := ""
	if len(locale) > 0 {
		targetLang = locale[0]
	}

	for i := range slice {
		lang := slice[i].Language
		text := slice[i].Text

		switch {
		case lang == targetLang:
			return text
		case lang == "en" && english == "":
			english = text
		case fallback == "":
			fallback = text
		}
	}

	if english != "" {
		return english
	}

	return fallback
}

// GetTexts processes multiple LocalizedTextSlice instances and returns their
// corresponding texts for a given locale.
//
// This function applies the same fallback mechanism as GetText to each slice
// in the input array. It is useful for batch processing of localized content.
//
// Example:
//
//	slices := []types.LocalizedTextSlice{
//		{{Text: "Hallo", Language: "de"}, {Text: "Hello", Language: "en"}},
//		{{Text: "Welt", Language: "de"}, {Text: "World", Language: "en"}},
//	}
//	texts := GetTexts(slices, "de") // Returns ["Hallo", "Welt"]
//
// Parameters:
//   - slices: Array of LocalizedTextSlice instances to process
//   - locale: Optional locale code for text selection
//
// Returns:
//   - []string: Array of texts in the requested locale
func GetTexts(slices []types.LocalizedTextSlice, locale ...string) []string {
	texts := make([]string, len(slices))
	for i, slice := range slices {
		texts[i] = GetText(slice, locale...)
	}

	return texts
}

// HasLanguage checks whether a LocalizedTextSlice contains text in the
// specified language.
//
// This function is useful for validating the availability of content in
// a particular language before attempting to retrieve it.
//
// Example:
//
//	texts := types.LocalizedTextSlice{
//		{Text: "Hello", Language: "en"},
//		{Text: "Hallo", Language: "de"},
//	}
//	hasGerman := HasLanguage(texts, "de") // Returns true
//	hasFrench := HasLanguage(texts, "fr") // Returns false
//
// Parameters:
//   - slice: The LocalizedTextSlice to check
//   - language: The language code to look for
//
// Returns:
//   - bool: True if the language is present, false otherwise
func HasLanguage(slice types.LocalizedTextSlice, language string) bool {
	for _, text := range slice {
		if text.Language == language {
			return true
		}
	}

	return false
}

// GetLanguages returns a list of all unique languages present in a
// LocalizedTextSlice.
//
// This function is useful for determining the available language options
// for a piece of content.
//
// Example:
//
//	texts := types.LocalizedTextSlice{
//		{Text: "Hello", Language: "en"},
//		{Text: "Hallo", Language: "de"},
//		{Text: "Hello", Language: "en"}, // Duplicate
//	}
//	languages := GetLanguages(texts) // Returns ["en", "de"]
//
// Parameters:
//   - slice: The LocalizedTextSlice to analyze
//
// Returns:
//   - []string: Array of unique language codes
func GetLanguages(slice types.LocalizedTextSlice) []string {
	languages := make(map[string]struct{})
	for _, text := range slice {
		languages[text.Language] = struct{}{}
	}

	result := make([]string, 0, len(languages))
	for lang := range languages {
		result = append(result, lang)
	}

	return result
}
