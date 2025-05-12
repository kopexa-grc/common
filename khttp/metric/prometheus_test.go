// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package metric

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetric(t *testing.T) {
	// setup router
	r := chi.NewRouter()
	r.Use(Middleware)

	r.Get("/clusters/{clusterId}", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/nodes/{nodeId}", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// do call
	var wg sync.WaitGroup

	for i := 0; i < 3000; i++ {
		wg.Add(1)

		go func() {
			req := httptest.NewRequest(http.MethodGet, "/clusters/foo", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			wg.Done()
		}()
	}

	for i := 0; i < 3000; i++ {
		wg.Add(1)

		go func() {
			req := httptest.NewRequest(http.MethodGet, "/nodes/foo", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			wg.Done()
		}()
	}

	wg.Wait()

	// check stats
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	Handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(data), "kopexa_http_time_to_first_header_duration_seconds_bucket")
	assert.Contains(t, string(data), "kopexa_http_response_size_bytes_bucket")
	assert.Contains(t, string(data), "kopexa_http_request_size_bytes_bucket")
	assert.Contains(t, string(data), "kopexa_http_request_duration_seconds_bucket")
	assert.Contains(t, string(data), "kopexa_http_in_flight_requests")
	assert.Contains(t, string(data), "kopexa_http_api_requests_total")
	assert.Contains(t, string(data), "kopexa_http_request_duration_seconds_count{code=\"200\",handler=\"/nodes/{nodeId}\",method=\"get\"} 3000")
	assert.Contains(t, string(data), "kopexa_http_response_size_bytes_count{code=\"200\",handler=\"/clusters/{clusterId}\",method=\"get\"} 3000")
	assert.Contains(t, string(data), "kopexa_http_request_size_bytes_count{code=\"200\",handler=\"/nodes/{nodeId}\",method=\"get\"} 3000")
}
