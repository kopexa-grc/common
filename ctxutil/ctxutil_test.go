// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package ctxutil

import (
	"context"
	"testing"
)

func TestWith(t *testing.T) {
	ctx := context.Background()
	value := "test-value"

	// Test storing a value
	newCtx := With(ctx, value)
	if newCtx == ctx {
		t.Error("With() should return a new context")
	}

	// Test retrieving the stored value
	if v, ok := From[string](newCtx); !ok || v != value {
		t.Errorf("From() = %v, %v; want %v, true", v, ok, value)
	}
}

func TestFrom(t *testing.T) {
	ctx := context.Background()
	value := "test-value"

	// Test retrieving non-existent value
	if v, ok := From[string](ctx); ok || v != "" {
		t.Errorf("From() = %v, %v; want \"\", false", v, ok)
	}

	// Test retrieving existing value
	ctx = With(ctx, value)
	if v, ok := From[string](ctx); !ok || v != value {
		t.Errorf("From() = %v, %v; want %v, true", v, ok, value)
	}
}

func TestMustFrom(t *testing.T) {
	ctx := context.Background()
	value := "test-value"

	// Test panic on non-existent value
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustFrom() should panic when value is not found")
		}
	}()
	MustFrom[string](ctx)

	// Test retrieving existing value
	ctx = With(ctx, value)
	if v := MustFrom[string](ctx); v != value {
		t.Errorf("MustFrom() = %v; want %v", v, value)
	}
}

func TestFromOr(t *testing.T) {
	ctx := context.Background()
	value := "test-value"
	defaultValue := "default-value"

	// Test with non-existent value
	if v := FromOr(ctx, defaultValue); v != defaultValue {
		t.Errorf("FromOr() = %v; want %v", v, defaultValue)
	}

	// Test with existing value
	ctx = With(ctx, value)
	if v := FromOr(ctx, defaultValue); v != value {
		t.Errorf("FromOr() = %v; want %v", v, value)
	}
}

func TestFromOrFunc(t *testing.T) {
	ctx := context.Background()
	value := "test-value"
	defaultValue := "default-value"

	// Test with non-existent value
	if v := FromOrFunc(ctx, func() string { return defaultValue }); v != defaultValue {
		t.Errorf("FromOrFunc() = %v; want %v", v, defaultValue)
	}

	// Test with existing value
	ctx = With(ctx, value)
	if v := FromOrFunc(ctx, func() string { return defaultValue }); v != value {
		t.Errorf("FromOrFunc() = %v; want %v", v, value)
	}
}

func TestTypeSafety(t *testing.T) {
	ctx := context.Background()
	stringValue := "string-value"
	intValue := 42

	// Store different types
	ctx = With(ctx, stringValue)
	ctx = With(ctx, intValue)

	// Test retrieving string
	if v, ok := From[string](ctx); !ok || v != stringValue {
		t.Errorf("From[string]() = %v, %v; want %v, true", v, ok, stringValue)
	}

	// Test retrieving int
	if v, ok := From[int](ctx); !ok || v != intValue {
		t.Errorf("From[int]() = %v, %v; want %v, true", v, ok, intValue)
	}

	// Test type mismatch
	if v, ok := From[int](ctx); ok && v == 0 {
		t.Error("From[int]() should not return zero value when type mismatch")
	}
}

func TestNestedContext(t *testing.T) {
	ctx := context.Background()
	value1 := "value1"
	value2 := "value2"

	// Create nested contexts
	ctx1 := With(ctx, value1)
	ctx2 := With(ctx1, value2)

	// Test retrieving from nested context
	if v, ok := From[string](ctx2); !ok || v != value2 {
		t.Errorf("From() = %v, %v; want %v, true", v, ok, value2)
	}

	// Test retrieving from parent context
	if v, ok := From[string](ctx1); !ok || v != value1 {
		t.Errorf("From() = %v, %v; want %v, true", v, ok, value1)
	}
}

func TestComplexStruct(t *testing.T) {
	type ComplexStruct struct {
		ID      int
		Name    string
		Tags    []string
		Details map[string]interface{}
	}

	ctx := context.Background()
	value := ComplexStruct{
		ID:   1,
		Name: "test",
		Tags: []string{"tag1", "tag2"},
		Details: map[string]interface{}{
			"key1": 42,
			"key2": "value",
		},
	}

	// Test storing complex struct
	ctx = With(ctx, value)

	// Test retrieving complex struct
	if v, ok := From[ComplexStruct](ctx); !ok {
		t.Error("From() should return true for complex struct")
	} else if v.ID != value.ID || v.Name != value.Name {
		t.Errorf("From() = %+v; want %+v", v, value)
	}
}

func TestNilValues(t *testing.T) {
	ctx := context.Background()

	// Test storing nil pointer
	var nilPtr *string
	ctx = With(ctx, nilPtr)

	// Test retrieving nil pointer
	if v, ok := From[*string](ctx); !ok {
		t.Error("From() should return true for nil pointer")
	} else if v != nil {
		t.Errorf("From() = %v; want nil", v)
	}
}

func TestNestedContextWithDifferentTypes(t *testing.T) {
	ctx := context.Background()

	// Create nested contexts with different types
	ctx1 := With(ctx, "string-value")
	ctx2 := With(ctx1, 42)
	ctx3 := With(ctx2, true)

	// Test retrieving all values from deepest context
	if v, ok := From[string](ctx3); !ok || v != "string-value" {
		t.Errorf("From[string]() = %v, %v; want %v, true", v, ok, "string-value")
	}
	if v, ok := From[int](ctx3); !ok || v != 42 {
		t.Errorf("From[int]() = %v, %v; want %v, true", v, ok, 42)
	}
	if v, ok := From[bool](ctx3); !ok || v != true {
		t.Errorf("From[bool]() = %v, %v; want %v, true", v, ok, true)
	}
}

func TestContextPerformance(t *testing.T) {
	ctx := context.Background()
	const depth = 1000

	// Create deeply nested context
	for i := 0; i < depth; i++ {
		ctx = With(ctx, i)
	}

	// Test retrieving value from deepest level
	if v, ok := From[int](ctx); !ok || v != depth-1 {
		t.Errorf("From[int]() = %v, %v; want %v, true", v, ok, depth-1)
	}

	// Test retrieving value from middle level
	middleCtx := ctx
	for i := 0; i < depth/2; i++ {
		middleCtx = context.WithValue(middleCtx, key[int]{}, i)
	}
	if v, ok := From[int](middleCtx); !ok || v != depth/2-1 {
		t.Errorf("From[int]() = %v, %v; want %v, true", v, ok, depth/2-1)
	}
}

func TestConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	const goroutines = 10
	const iterations = 1000

	// Create a context with initial value
	ctx = With(ctx, 0)

	// Test concurrent access
	done := make(chan bool)
	for i := 0; i < goroutines; i++ {
		go func() {
			for j := 0; j < iterations; j++ {
				// Read value
				if v, ok := From[int](ctx); !ok {
					t.Error("From() should return true")
				} else if v < 0 {
					t.Errorf("From() = %v; want >= 0", v)
				}
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < goroutines; i++ {
		<-done
	}
}
