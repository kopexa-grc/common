// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import "context"

// BatchCheckGate is consulted before every BatchCheckObjectAccess call.
//
// It may block until the call is allowed to run -- it must then give up when
// ctx is done and return ctx's error -- and it may lower how many batches the
// call sends to OpenFGA at once (maxParallel; 0 keeps the SDK default). The
// caller runs release once the call has finished, successful or not.
//
// The gate exists because a single bulk read can fan out into hundreds of
// batch requests within seconds, and OpenFGA does not survive that. Who gets
// throttled, and how hard, is the application's policy, so the client only
// provides the seam.
type BatchCheckGate func(ctx context.Context, checks []AccessCheck) (release func(), maxParallel int32, err error)

// WithBatchCheckGate installs a gate that every BatchCheckObjectAccess call
// passes before it reaches OpenFGA (see BatchCheckGate).
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithBatchCheckGate(myGate),
//	)
func WithBatchCheckGate(gate BatchCheckGate) Option {
	return func(c *Client) {
		c.batchCheckGate = gate
	}
}
