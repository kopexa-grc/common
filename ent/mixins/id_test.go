// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package mixins

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDisplayIDGeneration tests the generation of display IDs
func TestDisplayIDGeneration(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		length int
	}{
		{"6 character ID", "test-id-1", 6},
		{"8 character ID", "test-id-2", 8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			displayID := generateShortCharID(tc.input, tc.length)
			assert.Len(t, displayID, tc.length)
			for _, c := range displayID {
				assert.True(t, (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
			}
		})
	}
}

func TestCollisionResistance(t *testing.T) {
	const iterations = 10000
	seen := make(map[string]bool)
	for i := 0; i < iterations; i++ {
		id := "test-id-" + string(rune(i))
		displayID := generateShortCharID(id, 6)
		if seen[displayID] {
			t.Errorf("Collision detected for display ID: %s", displayID)
		}
		seen[displayID] = true
	}
}

func TestDisplayIDConsistency(t *testing.T) {
	input := "test-id-123"
	length := 6
	displayID1 := generateShortCharID(input, length)
	displayID2 := generateShortCharID(input, length)
	assert.Equal(t, displayID1, displayID2)
}
