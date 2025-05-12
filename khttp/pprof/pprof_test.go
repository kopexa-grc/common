// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package pprof

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Endpoints(t *testing.T) {
	h := Handler()

	endpoints := []string{
		PathPrefixPProf + "/",
		PathPrefixPProf + "/cmdline",
		PathPrefixPProf + "/profile",
		PathPrefixPProf + "/symbol",
		PathPrefixPProf + "/trace",
	}

	for _, ep := range endpoints {
		t.Run(ep, func(t *testing.T) {
			req := httptest.NewRequest("GET", ep, nil)
			rw := httptest.NewRecorder()
			h.ServeHTTP(rw, req)
			assert.NotEqual(t, http.StatusNotFound, rw.Code, "Endpoint %s should not return 404", ep)
			assert.NotEqual(t, http.StatusMethodNotAllowed, rw.Code, "Endpoint %s should not return 405", ep)
		})
	}
}
