// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package summarizer

import "errors"

var (
	// ErrSentenceEmpty is used to denote required sentences that needs to be summarized
	ErrSentenceEmpty = errors.New("you cannot summarize an empty string")
)
