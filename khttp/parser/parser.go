// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package parser

import (
	"net/http"
	"strconv"
	"time"
)

// ParseQueryInt parses an integer from a query parameter, returning a default value if parsing fails.
func ParseQueryInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	return parseIntOrDefault(value, defaultValue)
}

// parseIntOrDefault converts a string to an int, returning a default value if parsing fails.
func parseIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// ParseTimeParam parses a time parameter in RFC3339 format
func ParseTimeParam(r *http.Request, key string) *time.Time {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, param)
	if err != nil {
		return nil
	}

	return &t
}
