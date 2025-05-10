// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package sessions

// Default configuration values
const (
	// DefaultMaxAge is the default session duration in seconds (1 hour)
	DefaultMaxAge = 3600

	// DefaultSigningKey is the default signing key (32 bytes)
	DefaultSigningKey = "my-signing-secret"

	// DefaultEncryptionKey is the default encryption key (32 bytes)
	DefaultEncryptionKey = "encryptionsecret"
)

// Cookie configuration
const (
	// CookiePath is the default path for cookies
	CookiePath = "/"

	// CookieSameSiteLax is the default SameSite attribute
	CookieSameSiteLax = "Lax"

	// CookieSameSiteStrict is the strict SameSite attribute
	CookieSameSiteStrict = "Strict"

	// CookieSameSiteNone is the none SameSite attribute
	CookieSameSiteNone = "None"
)

// Key length validation
const (
	// MinKeyLength is the minimum allowed key length
	MinKeyLength = 16

	// DefaultKeyLength is the recommended key length
	DefaultKeyLength = 32

	// MaxKeyLength is the maximum allowed key length
	MaxKeyLength = 64

	// SessionIDLength is the length of the session ID in bytes (256 bits)
	SessionIDLength = 32
)
