package to

import (
	"testing"
	"time"
)

func TestBool(t *testing.T) {
	v := true
	if *Bool(v) != v {
		t.Errorf("Bool() = %v, want %v", *Bool(v), v)
	}
}

func TestString(t *testing.T) {
	v := "test"
	if *String(v) != v {
		t.Errorf("String() = %v, want %v", *String(v), v)
	}
}

func TestInt(t *testing.T) {
	v := 42
	if *Int(v) != v {
		t.Errorf("Int() = %v, want %v", *Int(v), v)
	}
}

func TestInt8(t *testing.T) {
	v := int8(42)
	if *Int8(v) != v {
		t.Errorf("Int8() = %v, want %v", *Int8(v), v)
	}
}

func TestInt16(t *testing.T) {
	v := int16(42)
	if *Int16(v) != v {
		t.Errorf("Int16() = %v, want %v", *Int16(v), v)
	}
}

func TestInt32(t *testing.T) {
	v := int32(42)
	if *Int32(v) != v {
		t.Errorf("Int32() = %v, want %v", *Int32(v), v)
	}
}

func TestInt64(t *testing.T) {
	v := int64(42)
	if *Int64(v) != v {
		t.Errorf("Int64() = %v, want %v", *Int64(v), v)
	}
}

func TestUint(t *testing.T) {
	v := uint(42)
	if *Uint(v) != v {
		t.Errorf("Uint() = %v, want %v", *Uint(v), v)
	}
}

func TestUint8(t *testing.T) {
	v := uint8(42)
	if *Uint8(v) != v {
		t.Errorf("Uint8() = %v, want %v", *Uint8(v), v)
	}
}

func TestUint16(t *testing.T) {
	v := uint16(42)
	if *Uint16(v) != v {
		t.Errorf("Uint16() = %v, want %v", *Uint16(v), v)
	}
}

func TestUint32(t *testing.T) {
	v := uint32(42)
	if *Uint32(v) != v {
		t.Errorf("Uint32() = %v, want %v", *Uint32(v), v)
	}
}

func TestUint64(t *testing.T) {
	v := uint64(42)
	if *Uint64(v) != v {
		t.Errorf("Uint64() = %v, want %v", *Uint64(v), v)
	}
}

func TestFloat32(t *testing.T) {
	v := float32(42.0)
	if *Float32(v) != v {
		t.Errorf("Float32() = %v, want %v", *Float32(v), v)
	}
}

func TestFloat64(t *testing.T) {
	v := float64(42.0)
	if *Float64(v) != v {
		t.Errorf("Float64() = %v, want %v", *Float64(v), v)
	}
}

func TestComplex64(t *testing.T) {
	v := complex64(42 + 0i)
	if *Complex64(v) != v {
		t.Errorf("Complex64() = %v, want %v", *Complex64(v), v)
	}
}

func TestComplex128(t *testing.T) {
	v := complex128(42 + 0i)
	if *Complex128(v) != v {
		t.Errorf("Complex128() = %v, want %v", *Complex128(v), v)
	}
}

func TestByte(t *testing.T) {
	v := byte(42)
	if *Byte(v) != v {
		t.Errorf("Byte() = %v, want %v", *Byte(v), v)
	}
}

func TestRune(t *testing.T) {
	v := rune('A')
	if *Rune(v) != v {
		t.Errorf("Rune() = %v, want %v", *Rune(v), v)
	}
}

func TestTime(t *testing.T) {
	v := time.Now()
	if *Time(v) != v {
		t.Errorf("Time() = %v, want %v", *Time(v), v)
	}
}

func TestDuration(t *testing.T) {
	v := time.Hour
	if *Duration(v) != v {
		t.Errorf("Duration() = %v, want %v", *Duration(v), v)
	}
}
