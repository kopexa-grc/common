// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package llm

import "errors"

// Common errors for the LLM package
var (
	ErrConfigRequired      = errors.New("config must not be nil")
	ErrUnsupportedProvider = errors.New("unsupported llm provider")
	ErrInvalidCredentials  = errors.New("invalid credentials provided")
)
