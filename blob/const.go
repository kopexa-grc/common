// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob

import (
	"time"

	kerr "github.com/kopexa-grc/common/errors"
)

var errClosed = kerr.Newf(kerr.FailedPrecondition, nil, "blob: Bucket has been closed")

// DefaultSignedURLExpiry is the default duration for SignedURLOptions.Expiry.
const DefaultSignedURLExpiry = 1 * time.Hour

const (
	hotAccessTier = "hot"
)

const (
	containerAccessType = "container"
	blobAccessType      = "blob"
	privateAccessType   = "private"
)

const (
	PublicContainer = "public"
)
