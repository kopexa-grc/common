// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package cors

// HTTP Methods
const (
	MethodHead    = "HEAD"
	MethodOptions = "OPTIONS"
	MethodGet     = "GET"
	MethodPost    = "POST"
	MethodPut     = "PUT"
	MethodPatch   = "PATCH"
	MethodDelete  = "DELETE"
)

// HTTP Headers
const (
	HeaderAccept                  = "Accept"
	HeaderAcceptLanguage          = "Accept-Language"
	HeaderAuthorization           = "Authorization"
	HeaderContentType             = "Content-Type"
	HeaderCSRFToken               = "X-CSRF-Token"
	HeaderContentDisposition      = "Content-Disposition"
	HeaderContentTransferEncoding = "Content-Transfer-Encoding"
	HeaderCredentials             = "Credentials"
	HeaderReferrerPolicy          = "Referrer-Policy"
	HeaderReferrerPolicyAlt       = "Referrerpolicy"

	// Exposed Headers
	HeaderLocation        = "Location"
	HeaderContentLocation = "Content-Location"
	HeaderDate            = "Date"
	HeaderETag            = "ETag"
)

// Time Constants
const (
	DefaultMaxAge = 60 * 60 * 6 // 6 hours in seconds
)
