// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package summarizer

import "errors"

// Common errors for the summarizer package
var (
	ErrConfigRequired    = errors.New("config must not be nil")
	ErrLLMConfigRequired = errors.New("LLM config is required for LLM summarization")
	ErrUnsupportedType   = errors.New("unsupported summarizer type")
)
