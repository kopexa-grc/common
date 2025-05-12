// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test cases
	tests := []struct {
		name           string
		method         string
		origin         string
		requestMethod  string
		expectedStatus int
		checkHeaders   func(t *testing.T, w *httptest.ResponseRecorder)
	}{
		{
			name:           "GET request with origin",
			method:         "GET",
			origin:         "http://localhost:8080",
			requestMethod:  "",
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "http://localhost:8080", w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
			},
		},
		{
			name:           "GET request without origin",
			method:         "GET",
			origin:         "",
			requestMethod:  "",
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
			},
		},
		{
			name:           "OPTIONS request with origin and method",
			method:         "OPTIONS",
			origin:         "http://localhost:8080",
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Equal(t, "http://localhost:8080", w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
				assert.Equal(t, "GET", w.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "21600", w.Header().Get("Access-Control-Max-Age"))
			},
		},
		{
			name:           "OPTIONS request without origin",
			method:         "OPTIONS",
			origin:         "",
			requestMethod:  "GET",
			expectedStatus: http.StatusOK,
			checkHeaders: func(t *testing.T, w *httptest.ResponseRecorder) {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"))
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Methods"))
				assert.Empty(t, w.Header().Get("Access-Control-Max-Age"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.requestMethod != "" {
				req.Header.Set("Access-Control-Request-Method", tt.requestMethod)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Apply middleware
			middleware := Middleware(handler)

			// Serve request
			middleware.ServeHTTP(w, req)

			// Check response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check headers
			tt.checkHeaders(t, w)
		})
	}
}

func TestAllowExtraHeader(t *testing.T) {
	// Save original headers
	originalHeaders := make([]string, len(Configuration.AllowedHeaders))
	copy(originalHeaders, Configuration.AllowedHeaders)

	// Test adding a new header
	newHeader := "X-Custom-Header"
	AllowExtraHeader(newHeader)

	// Check if header was added
	assert.Contains(t, Configuration.AllowedHeaders, newHeader)

	// Restore original headers
	Configuration.AllowedHeaders = originalHeaders
}

func TestExposeExtraHeader(t *testing.T) {
	// Save original headers
	originalHeaders := make([]string, len(Configuration.ExposedHeaders))
	copy(originalHeaders, Configuration.ExposedHeaders)

	// Test adding a new header
	newHeader := "X-Custom-Header"
	ExposeExtraHeader(newHeader)

	// Check if header was added
	assert.Contains(t, Configuration.ExposedHeaders, newHeader)

	// Restore original headers
	Configuration.ExposedHeaders = originalHeaders
}
