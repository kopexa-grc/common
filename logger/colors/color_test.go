// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package colors

import (
	"testing"

	"github.com/muesli/termenv"
	"github.com/stretchr/testify/assert"
)

func TestProfileName(t *testing.T) {
	tests := []struct {
		name     string
		profile  termenv.Profile
		expected string
	}{
		{
			name:     "Ascii profile",
			profile:  termenv.Ascii,
			expected: "Ascii",
		},
		{
			name:     "ANSI profile",
			profile:  termenv.ANSI,
			expected: "ANSI",
		},
		{
			name:     "ANSI256 profile",
			profile:  termenv.ANSI256,
			expected: "ANSI256",
		},
		{
			name:     "TrueColor profile",
			profile:  termenv.TrueColor,
			expected: "TrueColor",
		},
		{
			name:     "Unknown profile",
			profile:  termenv.Profile(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProfileName(tt.profile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTheme(t *testing.T) {
	// Test that DefaultColorTheme is properly initialized
	assert.NotNil(t, DefaultColorTheme)

	// Test that all color fields are set
	assert.NotNil(t, DefaultColorTheme.Primary)
	assert.NotNil(t, DefaultColorTheme.Secondary)
	assert.NotNil(t, DefaultColorTheme.Disabled)
	assert.NotNil(t, DefaultColorTheme.Error)
	assert.NotNil(t, DefaultColorTheme.Success)
	assert.NotNil(t, DefaultColorTheme.Critical)
	assert.NotNil(t, DefaultColorTheme.High)
	assert.NotNil(t, DefaultColorTheme.Medium)
	assert.NotNil(t, DefaultColorTheme.Low)
	assert.NotNil(t, DefaultColorTheme.Good)
	assert.NotNil(t, DefaultColorTheme.Unknown)
}
