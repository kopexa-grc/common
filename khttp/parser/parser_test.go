// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package parser

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseQueryInt(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		key          string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid integer",
			query:        "value=42",
			key:          "value",
			defaultValue: 0,
			expected:     42,
		},
		{
			name:         "invalid integer",
			query:        "value=invalid",
			key:          "value",
			defaultValue: 0,
			expected:     0,
		},
		{
			name:         "missing key",
			query:        "other=42",
			key:          "value",
			defaultValue: 0,
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{
					RawQuery: tt.query,
				},
			}
			result := ParseQueryInt(req, tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseTimeParam(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		key      string
		expected *time.Time
	}{
		{
			name:  "valid time",
			query: "time=2023-01-01T12:00:00Z",
			key:   "time",
			expected: func() *time.Time {
				t, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")
				return &t
			}(),
		},
		{
			name:     "invalid time",
			query:    "time=invalid",
			key:      "time",
			expected: nil,
		},
		{
			name:     "missing key",
			query:    "other=2023-01-01T12:00:00Z",
			key:      "time",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL: &url.URL{
					RawQuery: tt.query,
				},
			}
			result := ParseTimeParam(req, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}
