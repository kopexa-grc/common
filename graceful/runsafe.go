// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package graceful

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Recoverable interface {
	Start(context.Context) error
}

// RunSafe runs fn() and recovers panics, returning them as error
func RunSafe(ctx context.Context, name string, r Recoverable) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Any("panic", r).Str("component", name).Msg("Recovered from panic")
			err = fmt.Errorf("panic in %s: %v", name, r)
		}
	}()
	return r.Start(ctx)
}
