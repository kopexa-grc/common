// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package router_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kopexa-grc/common/khttp/cors"
	"github.com/kopexa-grc/common/khttp/router"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testServiceName = "router.test"

func TestNewObservabilityRouter(t *testing.T) {
	systemServiceName, serviceNameSet := os.LookupEnv("SERVICE_NAME")
	require.NoError(t, os.Setenv("SERVICE_NAME", testServiceName))
	defer func() {
		if serviceNameSet {
			_ = os.Setenv("SERVICE_NAME", systemServiceName)
		} else {
			_ = os.Unsetenv("SERVICE_NAME")
		}
	}()

	r := router.NewObservabilityRouter()

	t.Run("metrics", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

		r.ServeHTTP(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(data), "go_threads")
	})
	t.Run("debug", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/debug/pprof/cmdline", nil)

		r.ServeHTTP(rec, req)
		resp := rec.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(data), "router.test")
	},
	)
}

func TestRouterCors(t *testing.T) {
	type test struct {
		name    string
		subject chi.Router
	}

	customRegistry := router.NewMiddlewareRegistry()
	customRegistry.Register(cors.Middleware)

	tests := []test{
		{
			name:    "default router",
			subject: router.New(),
		},
		{
			name:    "option enabled router",
			subject: router.NewWithOptions(router.WithMiddlewareRegistry(customRegistry)),
		},
	}

	lo.ForEach(tests, func(test test, _ int) {
		t.Run(test.name, func(t *testing.T) {
			test.subject.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusConflict) })
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodOptions, "/health", nil)
			req.Header.Add("Origin", "some-origin")
			req.Header.Add("Access-Control-Request-Method", "POST")
			req.Header.Add("Access-Control-Request-Headers", "Content-Type")
			test.subject.ServeHTTP(rec, req)
			resp := rec.Result()
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			data, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.NotContains(t, string(data), "router.test OK")
			assert.Equal(t, "some-origin", resp.Header.Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "POST", resp.Header.Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Content-Type", resp.Header.Get("Access-Control-Allow-Headers"))
			assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, "21600", resp.Header.Get("Access-Control-Max-Age"))
		})
	})
}
