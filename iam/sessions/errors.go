// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package sessions

import "errors"

// Common errors that can occur during session operations
var (
	ErrSigningKeyTooShort         = errors.New("signing key must be at least 32 bytes")
	ErrEncryptionKeyTooShort      = errors.New("encryption key must be at least 32 bytes")
	ErrMaxAgeMustBePositive       = errors.New("max age must be positive")
	ErrInvalidSameSite            = errors.New("invalid SameSite value")
	ErrSecureRequired             = errors.New("secure must be true for production")
	ErrHTTPOnlyRequired           = errors.New("httpOnly must be true for security")
	ErrSameSiteNoneRequiresSecure = errors.New("SameSite=None requires Secure=true")
	ErrBucketNameRequired         = errors.New("bucket name is required")
	ErrServerURLRequired          = errors.New("server URL is required")
	ErrSaveFailed                 = errors.New("save error")
	ErrLoadFailed                 = errors.New("load error")
)
