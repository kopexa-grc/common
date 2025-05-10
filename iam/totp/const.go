// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package totp

// Default configuration values
const (
	// DefaultLength is the default length of TOTP codes
	DefaultLength = 6

	// DefaultRecoveryCodeCount is the default number of recovery codes
	DefaultRecoveryCodeCount = 16

	// DefaultRecoveryCodeLength is the default length of recovery codes
	DefaultRecoveryCodeLength = 8

	// CodePeriod is the validity period of a TOTP code in seconds
	CodePeriod = 30

	// KeyTTL is the expiration time for a key in NATS KV store
	KeyTTL = 30

	// OTPExpiration is the expiration time for an OTP code in minutes
	OTPExpiration = 5

	// NumericCode is a string of numbers for code generation
	NumericCode = "0123456789"

	// AlphanumericCode is a string of numbers and letters for code generation
	AlphanumericCode = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Base32SecretLength is the length of the base32 secret for TOTP
	Base32SecretLength = 20
)
