// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package gql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPreloadString(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		fieldName string
		want      string
	}{
		{
			name:      "empty prefix",
			prefix:    "",
			fieldName: "user",
			want:      "user",
		},
		{
			name:      "with prefix",
			prefix:    "organization",
			fieldName: "user",
			want:      "organization.user",
		},
		{
			name:      "nested prefix",
			prefix:    "organization.users",
			fieldName: "profile",
			want:      "organization.users.profile",
		},
		{
			name:      "empty name",
			prefix:    "organization",
			fieldName: "",
			want:      "organization.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPreloadString(tt.prefix, tt.fieldName)
			assert.Equal(t, tt.want, result)
		})
	}
}
