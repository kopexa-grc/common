// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package totp

import "context"

// Manager manages the protocol for SMS/Email 2FA codes and TOTP codes
type Manager interface {
	// TOTPQRString returns a URL string used for TOTP code generation
	TOTPQRString(u *User) (string, error)
	// TOTPDecryptedSecret decrypts a TOTP secret
	TOTPDecryptedSecret(secret string) (string, error)
	// TOTPSecret creates a TOTP secret for code generation
	TOTPSecret(u *User) (string, error)
	// OTPCode creates a random OTP code and hash
	OTPCode(address string, method DeliveryMethod) (code, hash string, err error)
	// ValidateOTP checks if a User email/sms delivered OTP code is valid
	ValidateOTP(code, hash string) error
	// ValidateTOTP checks if a User TOTP code is valid
	ValidateTOTP(ctx context.Context, user *User, code string) error
	// GenerateRecoveryCodes creates a set of recovery codes for a user
	GenerateRecoveryCodes() []string
}
