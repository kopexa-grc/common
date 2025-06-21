// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob

import (
	"fmt"

	"github.com/kopexa-grc/common/blob/driver"
	kerr "github.com/kopexa-grc/common/errors"
)

func wrapError(b driver.Bucket, err error, key string) error {
	if err == nil {
		return nil
	}

	msg := "blob"
	if key != "" {
		msg += fmt.Sprintf(" (key %q)", key)
	}
	code := kerr.Code(err)

	return kerr.New(code, msg)
}
