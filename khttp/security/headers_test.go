// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package security

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecurityHeaders(t *testing.T) {
	var dummySecurityHandler http.HandlerFunc = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("ISE"))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}

	r := Headers(dummySecurityHandler)

	t.Run("check security headers", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		r.ServeHTTP(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "ISE", string(data))
		assert.Equal(t, "max-age=31536000; includeSubDomains; preload", resp.Header.Get("Strict-Transport-Security"))
		assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
		assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	})
}
