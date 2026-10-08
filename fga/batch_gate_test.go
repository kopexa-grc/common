// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// batchCheckServer stands in for OpenFGA's batch-check endpoint: it allows
// every check and records how many requests it served and how many were in
// flight at the same time.
type batchCheckServer struct {
	served      atomic.Int32
	inFlight    atomic.Int32
	maxInFlight atomic.Int32
}

func (s *batchCheckServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := s.inFlight.Add(1)
	defer s.inFlight.Add(-1)

	for {
		highest := s.maxInFlight.Load()
		if now <= highest || s.maxInFlight.CompareAndSwap(highest, now) {
			break
		}
	}

	// Long enough that batches sent in parallel overlap.
	time.Sleep(20 * time.Millisecond)

	var body struct {
		Checks []struct {
			CorrelationID string `json:"correlation_id"`
		} `json:"checks"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := make(map[string]map[string]bool, len(body.Checks))
	for _, c := range body.Checks {
		result[c.CorrelationID] = map[string]bool{"allowed": true}
	}

	s.served.Add(1)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"result": result})
}

func newGateTestClient(t *testing.T, gate BatchCheckGate) (*Client, *batchCheckServer) {
	t.Helper()

	srv := &batchCheckServer{}
	ts := httptest.NewServer(srv)
	t.Cleanup(ts.Close)

	opts := []Option{WithStoreID(mockStoreID)}
	if gate != nil {
		opts = append(opts, WithBatchCheckGate(gate))
	}

	c, err := NewClient(ts.URL, opts...)
	require.NoError(t, err)

	return c, srv
}

// accessChecks returns n can_view checks; 200 of them make the SDK send four
// batches of 50.
func accessChecks(n int) []AccessCheck {
	checks := make([]AccessCheck, n)
	for i := range checks {
		checks[i] = AccessCheck{
			SubjectID:  "u1",
			ObjectID:   fmt.Sprintf("o%d", i),
			ObjectType: "asset",
			Relation:   "can_view",
		}
	}

	return checks
}

func TestBatchCheckObjectAccess_GateLimitsParallelBatches(t *testing.T) {
	var released atomic.Int32

	gate := func(_ context.Context, checks []AccessCheck) (func(), int32, error) {
		assert.Len(t, checks, 200)
		return func() { released.Add(1) }, 1, nil
	}

	c, srv := newGateTestClient(t, gate)

	allowed, err := c.BatchCheckObjectAccess(context.Background(), accessChecks(200))
	require.NoError(t, err)

	assert.Len(t, allowed, 200)
	assert.Equal(t, int32(4), srv.served.Load(), "200 checks go out as four batches")
	assert.Equal(t, int32(1), srv.maxInFlight.Load(), "maxParallel 1 sends the batches one after another")
	assert.Equal(t, int32(1), released.Load(), "release runs exactly once")
}

func TestBatchCheckObjectAccess_GateZeroKeepsSDKParallelism(t *testing.T) {
	gate := func(context.Context, []AccessCheck) (func(), int32, error) {
		return func() {}, 0, nil
	}

	c, srv := newGateTestClient(t, gate)

	_, err := c.BatchCheckObjectAccess(context.Background(), accessChecks(200))
	require.NoError(t, err)

	assert.Greater(t, srv.maxInFlight.Load(), int32(1), "without a limit the SDK sends batches in parallel")
}

func TestBatchCheckObjectAccess_GateErrorSendsNothing(t *testing.T) {
	gate := func(ctx context.Context, _ []AccessCheck) (func(), int32, error) {
		return nil, 0, ctx.Err()
	}

	c, srv := newGateTestClient(t, gate)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.BatchCheckObjectAccess(ctx, accessChecks(10))
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
	assert.Equal(t, int32(0), srv.served.Load())
}

func TestBatchCheckObjectAccess_GateReleasesWhenTheCallFails(t *testing.T) {
	var mu sync.Mutex

	released := false
	gate := func(context.Context, []AccessCheck) (func(), int32, error) {
		return func() {
			mu.Lock()
			released = true
			mu.Unlock()
		}, 0, nil
	}

	// A server that refuses every request.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "nope", http.StatusBadRequest)
	}))
	t.Cleanup(ts.Close)

	c, err := NewClient(ts.URL, WithStoreID(mockStoreID), WithBatchCheckGate(gate))
	require.NoError(t, err)

	_, err = c.BatchCheckObjectAccess(context.Background(), accessChecks(10))
	require.Error(t, err)

	mu.Lock()
	defer mu.Unlock()
	assert.True(t, released)
}
