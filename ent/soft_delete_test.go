// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package ent

import (
	"context"
	"testing"
)

// testContextKey is a custom type for context keys to avoid using string keys
type testContextKey string

const testKey testContextKey = "test_key"

func TestSoftDeleteContext(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func(context.Context) context.Context
		checkSkipFunc func(context.Context) bool
		checkIsFunc   func(context.Context) bool
		expectedSkip  bool
		expectedIs    bool
	}{
		{
			name: "default context",
			setupContext: func(ctx context.Context) context.Context {
				return ctx
			},
			checkSkipFunc: CheckSkipSoftDelete,
			checkIsFunc:   CheckIsSoftDelete,
			expectedSkip:  false,
			expectedIs:    false,
		},
		{
			name:          "skip soft delete",
			setupContext:  SkipSoftDelete,
			checkSkipFunc: CheckSkipSoftDelete,
			checkIsFunc:   CheckIsSoftDelete,
			expectedSkip:  true,
			expectedIs:    false,
		},
		{
			name:          "is soft delete",
			setupContext:  IsSoftDelete,
			checkSkipFunc: CheckSkipSoftDelete,
			checkIsFunc:   CheckIsSoftDelete,
			expectedSkip:  false,
			expectedIs:    true,
		},
		{
			name: "nested context - skip soft delete",
			setupContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, testKey, "some_value")
				return SkipSoftDelete(ctx)
			},
			checkSkipFunc: CheckSkipSoftDelete,
			checkIsFunc:   CheckIsSoftDelete,
			expectedSkip:  true,
			expectedIs:    false,
		},
		{
			name: "nested context - is soft delete",
			setupContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, testKey, "some_value")
				return IsSoftDelete(ctx)
			},
			checkSkipFunc: CheckSkipSoftDelete,
			checkIsFunc:   CheckIsSoftDelete,
			expectedSkip:  false,
			expectedIs:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = tt.setupContext(ctx)

			// Test skip soft delete
			if got := tt.checkSkipFunc(ctx); got != tt.expectedSkip {
				t.Errorf("CheckSkipSoftDelete() = %v, want %v", got, tt.expectedSkip)
			}

			// Test is soft delete
			if got := tt.checkIsFunc(ctx); got != tt.expectedIs {
				t.Errorf("CheckIsSoftDelete() = %v, want %v", got, tt.expectedIs)
			}
		})
	}
}

func TestSoftDeleteContextCombinations(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func(context.Context) context.Context
		checkSkipFunc func(context.Context) bool
		checkIsFunc   func(context.Context) bool
		expectedSkip  bool
		expectedIs    bool
	}{
		{
			name: "skip soft delete takes precedence",
			setupContext: func(ctx context.Context) context.Context {
				ctx = IsSoftDelete(ctx)
				return SkipSoftDelete(ctx)
			},
			checkSkipFunc: CheckSkipSoftDelete,
			checkIsFunc:   CheckIsSoftDelete,
			expectedSkip:  true,
			expectedIs:    true, // IsSoftDelete should still be true as it's set first
		},
		{
			name: "is soft delete after skip",
			setupContext: func(ctx context.Context) context.Context {
				ctx = SkipSoftDelete(ctx)
				return IsSoftDelete(ctx)
			},
			checkSkipFunc: CheckSkipSoftDelete,
			checkIsFunc:   CheckIsSoftDelete,
			expectedSkip:  true,
			expectedIs:    true, // IsSoftDelete should be true as it's set last
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = tt.setupContext(ctx)

			// Test skip soft delete
			if got := tt.checkSkipFunc(ctx); got != tt.expectedSkip {
				t.Errorf("CheckSkipSoftDelete() = %v, want %v", got, tt.expectedSkip)
			}

			// Test is soft delete
			if got := tt.checkIsFunc(ctx); got != tt.expectedIs {
				t.Errorf("CheckIsSoftDelete() = %v, want %v", got, tt.expectedIs)
			}
		})
	}
}
