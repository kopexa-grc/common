// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"context"

	"github.com/openfga/go-sdk/client"
)

// Expand returns an OpenFGA expand request builder
// See: https://openfga.dev/docs/interacting/relationship-queries#expand
func (c *Client) Expand(ctx context.Context) client.SdkClientExpandRequestInterface {
	return c.client.Expand(ctx)
}
