// Copyright by the openlane
// SPDX-License-Identifier: Apache-2.0

// Package ctxutil provides type-safe context value management using generics.
// It allows storing and retrieving values from a context without type assertions
// and with compile-time type safety.
//
// Example:
//
//	type User struct {
//		ID   string
//		Name string
//	}
//
//	// Store a user in the context
//	ctx = ctxutil.With(ctx, User{ID: "123", Name: "John"})
//
//	// Retrieve the user
//	if user, ok := ctxutil.From[User](ctx); ok {
//		fmt.Printf("User: %+v\n", user)
//	}
//
//	// Use a default value if not found
//	user := ctxutil.FromOr(ctx, User{ID: "default"})
//
//	// Use a function to compute default value
//	user = ctxutil.FromOrFunc(ctx, func() User {
//		return User{ID: "computed"}
//	})
package ctxutil

import (
	"context"
)

// key is a private type used as a key in the context.
// Using a generic type ensures type safety and prevents key collisions.
type key[T any] struct{}

// With stores a value of type T in the context.
// The value can be retrieved using From[T] with the resulting context.
// The function uses a generic key type to ensure type safety and prevent key collisions.
//
// Example:
//
//	ctx = ctxutil.With(ctx, "my-value")
func With[T any](ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, key[T]{}, v)
}

// From retrieves a value of type T from the context.
// Returns the value and true if found, or the zero value and false if not found.
//
// Example:
//
//	if value, ok := ctxutil.From[string](ctx); ok {
//		fmt.Println(value)
//	}
func From[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(key[T]{}).(T)
	return v, ok
}

// MustFrom retrieves a value of type T from the context.
// Panics if the value is not found or if the type assertion fails.
// Use this function only when you are certain the value exists.
//
// Example:
//
//	value := ctxutil.MustFrom[string](ctx)
func MustFrom[T any](ctx context.Context) T {
	return ctx.Value(key[T]{}).(T)
}

// FromOr retrieves a value of type T from the context.
// Returns the default value if the value is not found.
// This is useful when you want to ensure a value is always returned.
//
// Example:
//
//	value := ctxutil.FromOr(ctx, "default-value")
func FromOr[T any](ctx context.Context, def T) T {
	v, ok := From[T](ctx)
	if !ok {
		return def
	}

	return v
}

// FromOrFunc retrieves a value of type T from the context.
// If the value is not found, calls the provided function to compute a default value.
// This is useful when the default value is expensive to compute or depends on runtime conditions.
//
// Example:
//
//	value := ctxutil.FromOrFunc(ctx, func() string {
//		return computeExpensiveDefault()
//	})
func FromOrFunc[T any](ctx context.Context, f func() T) T {
	v, ok := From[T](ctx)
	if !ok {
		return f()
	}

	return v
}
