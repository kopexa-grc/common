// Package ptr
//
// Overview
// The ptr package provides small, generic helpers for pointer ergonomics:
//   - To(v)        obtain *T for any value v (generic replacement for to.Int, to.String, ...)
//   - Deref(p, d)  nil-safe dereference returning d when p is nil
//   - Equal(a, b)  nil-safe pointer and value equality
//   - AllPtrFieldsNil(obj)  report whether all pointer fields in (pointer to) struct are nil
//
// Rationale
// Prior to Go generics, libraries often included a combinatorial explosion of helpers:
//   Bool, Int, Int64, String, Time, ...
// Generics permit a single implementation (To) instead. Concentrating these utilities here lets
// us deprecate the older to package and reduce API surface area.
//
// Migration From package to
//   to.String("x")  -> ptr.To("x")
//   to.Int64(5)     -> ptr.To(int64(5))
//   to.BoolValue(p) -> ptr.Deref(p, false)
//   to.TimeValue(p) -> ptr.Deref(p, time.Time{})
//
// Example
//  name := ptr.To("alice")
//  timeout := ptr.To(30 * time.Second)
//  enabled := ptr.Deref(nil, true) // true (default)
//
// License
// This package is distributed under the Business Source License 1.1 (BUSL-1.1).
