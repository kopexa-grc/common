// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package ptr

import "testing"

func TestTo(t *testing.T) {
	p1 := To(42)
	p2 := To(42)

	if p1 == p2 { // distinct allocations expected
		t.Fatalf("expected distinct pointers for identical values")
	}

	if *p1 != 42 || *p2 != 42 {
		t.Fatalf("unexpected values: %d %d", *p1, *p2)
	}
}

func TestDeref(t *testing.T) {
	var ip *int
	if v := Deref(ip, 7); v != 7 {
		t.Fatalf("expected default 7, got %d", v)
	}

	x := 9
	if v := Deref(&x, 7); v != 9 {
		t.Fatalf("expected 9, got %d", v)
	}
}

func TestEqual(t *testing.T) {
	if !Equal[int](nil, nil) {
		t.Fatal("nil,nil should be equal")
	}

	x := 3
	if Equal(&x, nil) || Equal(nil, &x) {
		t.Fatal("one nil should not be equal")
	}

	y := 3
	if !Equal(&x, &y) {
		t.Fatal("expected equal values")
	}

	z := 4
	if Equal(&x, &z) {
		t.Fatal("expected inequality for different values")
	}
}

type testStructAllNil struct {
	A *int
	B *string
}

type testStructMixed struct {
	A *int
	B *string
}

func TestAllPtrFieldsNil(t *testing.T) {
	if !AllPtrFieldsNil(&testStructAllNil{}) {
		t.Fatal("expected true for all nil pointer fields (pointer receiver)")
	}

	if !AllPtrFieldsNil(testStructAllNil{}) {
		t.Fatal("expected true for all nil pointer fields (value receiver)")
	}

	s := testStructMixed{}
	v := 1

	s.A = &v
	if AllPtrFieldsNil(s) {
		t.Fatal("expected false when a pointer field is non-nil")
	}

	var nilPtr *testStructAllNil
	if !AllPtrFieldsNil(nilPtr) { // typed nil pointer
		t.Fatal("expected true for typed nil pointer")
	}
	// Panic cases: non-struct kind
	assertPanic(t, func() { AllPtrFieldsNil(123) })
}

func assertPanic(t *testing.T, fn func()) {
	t.Helper()

	didPanic := false

	defer func() {
		if r := recover(); r != nil {
			didPanic = true
		}

		if !didPanic {
			t.Fatalf("expected panic, did not panic")
		}
	}()
	fn()
}
