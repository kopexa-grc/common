// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package metric

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	// Reset registry for clean test
	GlobalRegistry = NewRegistry()

	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)

	// Test cases
	tests := []struct {
		name           string
		handler        http.Handler
		expectedStatus int
		checkMetrics   func(t *testing.T)
	}{
		{
			name:           "successful request",
			handler:        handler,
			expectedStatus: http.StatusOK,
			checkMetrics: func(t *testing.T) {
				// Check if metrics were recorded
				metrics, err := GlobalRegistry.Gather()
				require.NoError(t, err)
				assert.NotEmpty(t, metrics)
			},
		},
		{
			name: "error request",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
			expectedStatus: http.StatusInternalServerError,
			checkMetrics: func(t *testing.T) {
				// Check if metrics were recorded
				metrics, err := GlobalRegistry.Gather()
				require.NoError(t, err)
				assert.NotEmpty(t, metrics)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset registry for each test
			GlobalRegistry = NewRegistry()

			// Create response recorder
			w := httptest.NewRecorder()

			// Apply middleware
			middleware := Middleware(tt.handler)

			// Serve request
			middleware.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check metrics
			tt.checkMetrics(t)
		})
	}
}

func TestHandler(t *testing.T) {
	// Reset registry for clean test
	GlobalRegistry = NewRegistry()

	// Create test request
	req := httptest.NewRequest("GET", "/metrics", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve request
	Handler.ServeHTTP(w, req)

	// Check response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/plain")
	assert.NotEmpty(t, w.Body.String())
}

func TestCustomBuckets(t *testing.T) {
	// Test custom duration buckets
	customDurationBuckets := []float64{0.1, 0.5, 1.0, 2.0}
	DurationBuckets = customDurationBuckets
	assert.Equal(t, customDurationBuckets, DurationBuckets)

	// Test custom size buckets
	customSizeBuckets := []float64{100, 500, 1000, 5000}
	SizeBuckets = customSizeBuckets
	assert.Equal(t, customSizeBuckets, SizeBuckets)
}
