// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package sessions

import (
	"context"

	"github.com/kopexa-grc/common/ctxutil"
)

// SessionContextKey is a type-safe key for storing sessions in the context
type SessionContextKey[T any] struct{}

// WithSession stores a session in the context
func WithSession[T any](ctx context.Context, session *Session[T]) context.Context {
	return ctxutil.With(ctx, session)
}

// FromSession retrieves a session from the context
func FromSession[T any](ctx context.Context) (*Session[T], bool) {
	return ctxutil.From[*Session[T]](ctx)
}

// MustFromSession retrieves a session from the context and panics if not found
func MustFromSession[T any](ctx context.Context) *Session[T] {
	return ctxutil.MustFrom[*Session[T]](ctx)
}

// FromSessionOr retrieves a session from the context or returns a default value
func FromSessionOr[T any](ctx context.Context, def *Session[T]) *Session[T] {
	return ctxutil.FromOr(ctx, def)
}

// FromSessionOrFunc retrieves a session from the context or computes a default value
func FromSessionOrFunc[T any](ctx context.Context, f func() *Session[T]) *Session[T] {
	return ctxutil.FromOrFunc(ctx, f)
}
