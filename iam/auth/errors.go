// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package auth

import "github.com/kopexa-grc/common/errors"

var (
	ErrNoAuthUser         = errors.NewNotFound("no auth user")
	ErrInvalidCredentials = errors.NewBadRequest("the provided credentials are missing or invalid")
)
