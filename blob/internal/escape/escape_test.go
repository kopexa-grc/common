// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package escape_test

import (
	"strings"
	"testing"

	"github.com/kopexa-grc/common/blob/internal/escape"
)

// WeirdStrings are unusual/weird strings for use in testing escaping.
// The keys are descriptive strings, the values are the weird strings.
var WeirdStrings = map[string]string{
	"fwdslashes":          "foo/bar/baz",
	"repeatedfwdslashes":  "foo//bar///baz",
	"dotdotslash":         "../foo/../bar/../../baz../",
	"backslashes":         "foo\\bar\\baz",
	"repeatedbackslashes": "..\\foo\\\\bar\\\\\\baz",
	"dotdotbackslash":     "..\\foo\\..\\bar\\..\\..\\baz..\\",
	"quote":               "foo\"bar\"baz",
	"spaces":              "foo bar baz",
	"startwithdigit":      "12345",
	"unicode":             strings.Repeat("☺", 3),
	// The ASCII characters 0-128, split up to avoid the possibly-escaped
	// versions from getting too long.
	"ascii-1": makeASCIIString(0, 16),
	"ascii-2": makeASCIIString(16, 32),
	"ascii-3": makeASCIIString(32, 48),
	"ascii-4": makeASCIIString(48, 64),
	"ascii-5": makeASCIIString(64, 80),
	"ascii-6": makeASCIIString(80, 96),
	"ascii-7": makeASCIIString(96, 112),
	"ascii-8": makeASCIIString(112, 128),
}

func makeASCIIString(start, end int) string {
	var s []byte
	for i := start; i < end; i++ {
		if i >= 'a' && i <= 'z' {
			continue
		}
		if i >= 'A' && i <= 'Z' {
			continue
		}
		if i >= '0' && i <= '9' {
			continue
		}
		s = append(s, byte(i))
	}
	return string(s)
}

func TestHexEscape(t *testing.T) {
	always := func([]rune, int) bool { return true }

	for _, tc := range []struct {
		description, s, want string
		should               func([]rune, int) bool
	}{
		{
			description: "empty string",
			s:           "",
			want:        "",
			should:      always,
		},
		{
			description: "first rune",
			s:           "hello world",
			want:        "__0x68__ello world",
			should:      func(_ []rune, i int) bool { return i == 0 },
		},
		{
			description: "last rune",
			s:           "hello world",
			want:        "hello worl__0x64__",
			should:      func(r []rune, i int) bool { return i == len(r)-1 },
		},
		{
			description: "runes in middle",
			s:           "hello  world",
			want:        "hello__0x20____0x20__world",
			should:      func(r []rune, i int) bool { return r[i] == ' ' },
		},
		{
			description: "unicode",
			s:           "☺☺",
			should:      always,
			want:        "__0x263a____0x263a__",
		},
	} {
		got := escape.HexEscape(tc.s, tc.should)
		if got != tc.want {
			t.Errorf("%s: got escaped %q want %q", tc.description, got, tc.want)
		}
		got = escape.HexUnescape(got)
		if got != tc.s {
			t.Errorf("%s: got unescaped %q want %q", tc.description, got, tc.s)
		}
	}
}

func TestHexEscapeUnescapeWeirdStrings(t *testing.T) {
	for name, s := range WeirdStrings {
		escaped := escape.HexEscape(s, func(r []rune, i int) bool { return !escape.IsASCIIAlphanumeric(r[i]) })
		unescaped := escape.HexUnescape(escaped)
		if unescaped != s {
			t.Errorf("%s: got unescaped %q want %q", name, unescaped, s)
		}
	}
}

func TestHexUnescapeOnInvalid(t *testing.T) {
	// Unescaping of valid escape sequences is tested in TestEscape.
	// This only tests invalid escape sequences, so Unescape is expected
	// to do nothing.
	for _, s := range []string{
		"0x68",
		"_0x68_",
		"__0x68_",
		"_0x68__",
		"__1x68__",
		"__0y68__",
		"__0xag__",       // invalid hex digit
		"__0x8fffffff__", // out of int32 range
	} {
		got := escape.HexUnescape(s)
		if got != s {
			t.Errorf("%s: got %q want %q", s, got, s)
		}
	}
}
