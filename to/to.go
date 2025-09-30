// Copyright 2023 Kopexa GRC. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.

// Package to provides utility functions for pointer management.
// It helps convert between pointers and values of various types.
package to

import "time"

// Ptr returns a pointer to the provided value.
//
// Deprecated: use ptr.To(v) from package ptr.
func Ptr[T any](v T) *T {
	return &v
}

// Bool returns a pointer to the given bool value.
//
// Deprecated: use ptr.To(v).
func Bool(v bool) *bool {
	return &v
}

// BoolValue returns the bool value from a pointer, or false if the pointer is nil.
//
// Deprecated: use ptr.Deref(v, false).
func BoolValue(v *bool) bool {
	if v == nil {
		return false
	}

	return *v
}

// String returns a pointer to the given string value.
//
// Deprecated: use ptr.To(v).
func String(v string) *string {
	return &v
}

// StringValue returns the string value from a pointer, or an empty string if the pointer is nil.
//
// Deprecated: use ptr.Deref(v, "").
func StringValue(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

// Int returns a pointer to the given int value.
//
// Deprecated: use ptr.To(v).
func Int(v int) *int {
	return &v
}

// Int8 returns a pointer to the given int8 value.
//
// Deprecated: use ptr.To(v).
func Int8(v int8) *int8 {
	return &v
}

// Int16 returns a pointer to the given int16 value.
//
// Deprecated: use ptr.To(v).
func Int16(v int16) *int16 {
	return &v
}

// Int32 returns a pointer to the given int32 value.
//
// Deprecated: use ptr.To(v).
func Int32(v int32) *int32 {
	return &v
}

// Int64 returns a pointer to the given int64 value.
//
// Deprecated: use ptr.To(v).
func Int64(v int64) *int64 {
	return &v
}

// Uint returns a pointer to the given uint value.
//
// Deprecated: use ptr.To(v).
func Uint(v uint) *uint {
	return &v
}

// Uint8 returns a pointer to the given uint8 value.
//
// Deprecated: use ptr.To(v).
func Uint8(v uint8) *uint8 {
	return &v
}

// Uint16 returns a pointer to the given uint16 value.
//
// Deprecated: use ptr.To(v).
func Uint16(v uint16) *uint16 {
	return &v
}

// Uint32 returns a pointer to the given uint32 value.
//
// Deprecated: use ptr.To(v).
func Uint32(v uint32) *uint32 {
	return &v
}

// Uint64 returns a pointer to the given uint64 value.
//
// Deprecated: use ptr.To(v).
func Uint64(v uint64) *uint64 {
	return &v
}

// Float32 returns a pointer to the given float32 value.
//
// Deprecated: use ptr.To(v).
func Float32(v float32) *float32 {
	return &v
}

// Float64 returns a pointer to the given float64 value.
//
// Deprecated: use ptr.To(v).
func Float64(v float64) *float64 {
	return &v
}

// Complex64 returns a pointer to the given complex64 value.
//
// Deprecated: use ptr.To(v).
func Complex64(v complex64) *complex64 {
	return &v
}

// Complex128 returns a pointer to the given complex128 value.
//
// Deprecated: use ptr.To(v).
func Complex128(v complex128) *complex128 {
	return &v
}

// Byte returns a pointer to the given byte value.
//
// Deprecated: use ptr.To(v).
func Byte(v byte) *byte {
	return &v
}

// Rune returns a pointer to the given rune value.
//
// Deprecated: use ptr.To(v).
func Rune(v rune) *rune {
	return &v
}

// Time returns a pointer to the given time.Time value.
//
// Deprecated: use ptr.To(v).
func Time(v time.Time) *time.Time {
	return &v
}

// TimeValue returns the time.Time value from a pointer, or the zero value if the pointer is nil.
//
// Deprecated: use ptr.Deref(v, time.Time{}).
func TimeValue(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}

	return *v
}

// Duration returns a pointer to the given time.Duration value.
//
// Deprecated: use ptr.To(v).
func Duration(v time.Duration) *time.Duration {
	return &v
}
