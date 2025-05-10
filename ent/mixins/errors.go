// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package mixins

import "errors"

// ErrUnexpectedMutationType is returned when an unexpected mutation type is encountered
// during audit logging.
var ErrUnexpectedMutationType = errors.New("unexpected mutation type for audit logging")
