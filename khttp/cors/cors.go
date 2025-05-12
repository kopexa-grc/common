// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package cors

import (
	"net/http"

	"github.com/go-chi/cors"
)

// Package cors contains the default CORS (Cross-Origin Resource Sharing)
// configuration for Kopexa APIs.
// In order change the default (extend it), it is highly recommended,
// to use the provided helper functions before creating the router.

// Configuration represents the default configuration
// used for Kopexa APIs
var Configuration = cors.Options{
	AllowOriginFunc: func(_ *http.Request, _ string) bool { return true },
	AllowedMethods: []string{
		MethodHead,
		MethodOptions,
		MethodGet,
		MethodPost,
		MethodPut,
		MethodPatch,
		MethodDelete,
	},
	AllowedHeaders: []string{
		HeaderAccept,
		HeaderAcceptLanguage,
		HeaderAuthorization,
		HeaderContentType,
		HeaderCSRFToken,
		HeaderContentDisposition,
		HeaderContentTransferEncoding,
		HeaderCredentials,
		HeaderReferrerPolicyAlt,
		HeaderReferrerPolicy,
	},
	ExposedHeaders: []string{
		// REST
		HeaderLocation,
		HeaderContentLocation,
		HeaderDate,
		HeaderETag,
	},
	AllowCredentials: true,
	MaxAge:           DefaultMaxAge,
}

// AllowExtraHeader allows an additional header to be passed
// to the API from a browser. Only use this function if this
// is specific to the service, otherwise extend this package.
func AllowExtraHeader(header string) {
	Configuration.AllowedHeaders = append(Configuration.AllowedHeaders, header)
}

// ExposeExtraHeader allows an additional header to be returned
// from the API to a browser. Only use this function if this
// is specific to the service, otherwise extend this package.
func ExposeExtraHeader(header string) {
	Configuration.ExposedHeaders = append(Configuration.ExposedHeaders, header)
}

// Middleware creates a cors middleware using the
// cors.Configuration
func Middleware(next http.Handler) http.Handler {
	return cors.New(Configuration).Handler(next)
}
