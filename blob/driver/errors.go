// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package driver

import "errors"

var (
	ErrUnsupportedMethod = errors.New("unsupported method")
	ErrCopyFailed        = errors.New("copy failed")
)
