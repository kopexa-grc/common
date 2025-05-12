// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package security

import "net/http"

// Headers defense in depth mechanisms to protect the malicious behavior on client side,
// and also it helps externals to determine the security posture of the enterprise services.
func Headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		w.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(w, r)
	})
}
