// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package totp

import "errors"

// Common errors that can occur during TOTP operations
var (
	ErrCannotGenerateRandomString     = errors.New("cannot generate random string")
	ErrCannotHashOTPString            = errors.New("cannot hash OTP string")
	ErrFailedToGenerateSecret         = errors.New("failed to generate TOTP secret")
	ErrCannotEncryptSecret            = errors.New("cannot encrypt secret")
	ErrFailedToGetSecretForQR         = errors.New("failed to get secret for QR code")
	ErrInvalidCode                    = errors.New("invalid code")
	ErrCannotDecryptSecret            = errors.New("cannot decrypt secret")
	ErrIncorrectCodeProvided          = errors.New("incorrect code provided")
	ErrCodeIsNoLongerValid            = errors.New("code is no longer valid")
	ErrFailedToValidateCode           = errors.New("failed to validate code")
	ErrNoSecretKey                    = errors.New("no secret key available")
	ErrNoSecretKeyForVersion          = errors.New("no secret key available for version")
	ErrFailedToCreateCipherBlock      = errors.New("failed to create cipher block")
	ErrFailedToCreateCipherText       = errors.New("failed to create cipher text")
	ErrFailedToDetermineSecretVersion = errors.New("failed to determine secret version")
	ErrCannotDecodeSecret             = errors.New("cannot decode secret")
	ErrCipherTextTooShort             = errors.New("cipher text too short")
	ErrFailedToHashCode               = errors.New("failed to hash code")
	ErrCannotDecodeOTPHash            = errors.New("cannot decode OTP hash")
	ErrInvalidOTPHashFormat           = errors.New("invalid OTP hash format")
	ErrNilJetStream                   = errors.New("nil JetStream context")
)
