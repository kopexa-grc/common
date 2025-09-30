// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package ptr provides utilities for working with pointers in a concise, generic way.
//
// Motivation
// The legacy package "to" contained a large set of small helper functions (Bool, Int64, String, ...)
// to obtain pointers to literals. With the advent of Go generics, a single generic helper (To)
// covers the majority of those use‑cases. This package offers a focused, documented replacement
// and additional helpers for safe dereferencing, comparison, and structural inspection.
//
// Migration
// All helpers in package "to" are now deprecated in favor of the following functions:
//
//	to.Bool(v)     -> ptr.To(v)
//	to.String(v)   -> ptr.To(v)
//	to.Int64(v)    -> ptr.To(v)
//	... etc.
//
// Value extraction helpers (e.g. to.BoolValue) can be replaced by ptr.Deref(ptrValue, zeroValue).
//
// Overview
//
//	To(v)                  - obtain *T for any value v (generic)
//	Deref(p, def)          - safely dereference *T returning def when p is nil
//	Equal(a, b)            - nil-safe pointer equality (value comparison when both non-nil)
//	AllPtrFieldsNil(obj)   - report whether all pointer fields of a (pointer to) struct are nil
//
// The functions are intentionally small and allocation-free beyond required pointer creation.
package ptr

import (
	"fmt"
	"reflect"
)

// To returns a pointer to the given value v.
// Each call allocates a new variable holding the value (identical to &v pattern).
// Prefer this over numerous type-specific helpers.
func To[T any](v T) *T {
	return &v
}

// Deref returns the value pointed to by ptr, or def when ptr is nil.
// This is a generic replacement for legacy *Value helpers.
func Deref[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}

	return def
}

// Equal reports whether two pointers are both nil, or both non-nil and their
// dereferenced values are equal. It never panics.
func Equal[T comparable](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}

	if a == nil {
		return true
	}

	return *a == *b
}

// AllPtrFieldsNil reports whether every pointer-typed field in a struct value is nil.
// It accepts either a struct or a pointer to a struct. A typed nil struct pointer returns true.
// Passing any other kind (slice, map, primitive, interface containing non-struct) causes a panic.
//
// Typical usage: detect whether an optional request payload (with pointer fields for partial
// updates) carries no user-provided values.
func AllPtrFieldsNil(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		panic(fmt.Sprintf("reflect.ValueOf() produced a non-valid Value for %#v", obj))
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}

		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && !v.Field(i).IsNil() {
			return false
		}
	}

	return true
}
