// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package tokens

import "github.com/kopexa-grc/common/errors"

// Common error definitions for token operations.
var (
	// ErrInviteTokenMissingEmail is returned when an organization invite token is created without an email.
	ErrInviteTokenMissingEmail = errors.NewBadRequest("email is required")

	// ErrExpirationIsRequired is returned when a token is created without an expiration time.
	ErrExpirationIsRequired = errors.NewBadRequest("expiration is required")

	// ErrFailedSigning is returned when token signing operations fail.
	ErrFailedSigning = errors.NewUnexpectedFailure("failed to sign token")

	// ErrTokenInvalid is returned when a token's signature cannot be verified.
	ErrTokenInvalid = errors.NewBadRequest("token is invalid")

	// ErrTokenExpired is returned when a token has passed its expiration time.
	ErrTokenExpired = errors.NewFailedPrecondition("token is expired")

	// ErrInvalidSecret is returned when the provided secret does not match the expected format.
	ErrInvalidSecret = errors.NewBadRequest("invalid secret")

	// ErrMissingEmail is returned when the token is attempted to be verified but the email is missing
	ErrMissingEmail = errors.New(errors.InvalidArgument, "unable to create verification token, email is missing")
	// ErrTokenMissingEmail is returned when the verification is missing an email address
	ErrTokenMissingEmail = errors.New(errors.InvalidArgument, "email verification token is missing email address")
)

// Cryptographic constants for token operations.
const (
	// nonceLength defines the length of the nonce used in token signing.
	nonceLength = 64

	// inviteExpirationDays defines the number of days an organization invite token remains valid.
	inviteExpirationDays = 14

	// keyLength defines the length of the HMAC key used in token signing.
	keyLength = 64

	expirationDays              = 7
	resetTokenExpirationMinutes = 15
)
