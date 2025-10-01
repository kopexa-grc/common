// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package tokens

import "errors"

var (
	// ErrTokenMissingUserID is returned during verification or validation when the
	// ResetToken struct lacks a UserID (logical integrity failure).
	ErrTokenMissingUserID = errors.New("reset token is missing user id")
	// ErrMissingUserID is returned at construction time (NewResetToken) when the
	// caller supplies an empty user id.
	ErrMissingUserID = errors.New("unable to create reset token, user id is required")
)
