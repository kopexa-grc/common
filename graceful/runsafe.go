// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package graceful

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

var ErrPanicInFunction = errors.New("panic in function")

type Recoverable interface {
	Start(context.Context) error
}

// RunSafe runs fn() and recovers panics, returning them as error
func RunSafe(ctx context.Context, name string, r Recoverable) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Any(LogFieldPanic, r).Str(LogFieldComponent, name).Msg(LogMsgRecoveredFromPanic)
			err = fmt.Errorf("%w: %s: %+v", ErrPanicInFunction, name, r)
		}
	}()

	return r.Start(ctx)
}
