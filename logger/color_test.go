// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package logger

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// mockTerminal is a mock for os.File that simulates a terminal.
type mockTerminal struct {
	bytes.Buffer
}

func (m *mockTerminal) Fd() uintptr { return os.Stdout.Fd() }
func (m *mockTerminal) Stat() (os.FileInfo, error) {
	return &mockFileInfo{mode: os.ModeCharDevice}, nil
}

// mockFileInfo implements os.FileInfo for testing.
type mockFileInfo struct{ mode os.FileMode }

func (m *mockFileInfo) Name() string       { return "mock" }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return false }
func (m *mockFileInfo) Sys() interface{}   { return nil }

func TestNewColorWriter(t *testing.T) {
	tests := []struct {
		name     string
		out      interface{}
		expected *ColorWriter
	}{
		{
			name: "nil output",
			out:  nil,
			expected: &ColorWriter{
				out:    os.Stdout,
				colors: make(map[string]string),
			},
		},
		{
			name: "bytes.Buffer output",
			out:  &bytes.Buffer{},
			expected: &ColorWriter{
				out:    &bytes.Buffer{},
				colors: make(map[string]string),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out io.Writer
			if tt.out != nil {
				out = tt.out.(io.Writer)
			}

			cw := NewColorWriter(out)
			assert.NotNil(t, cw)
			assert.NotEmpty(t, cw.colors)
		})
	}
}

func TestColorWriter_Write(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		writer   io.Writer
		expected string
	}{
		{
			name:     "non-terminal output",
			input:    `{"level":"info","message":"test"}`,
			writer:   &bytes.Buffer{},
			expected: `{"level":"info","message":"test"}`,
		},
		{
			name:     "terminal output with color",
			input:    `{"level":"info","message":"test"}`,
			writer:   &mockTerminal{},
			expected: `{"level":"info","message":"test"}`,
		},
		{
			name:     "empty input",
			input:    "",
			writer:   &bytes.Buffer{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw := NewColorWriter(tt.writer)
			n, err := cw.Write([]byte(tt.input))
			assert.NoError(t, err)
			assert.Equal(t, len(tt.input), n)

			if buf, ok := tt.writer.(*bytes.Buffer); ok {
				assert.Contains(t, buf.String(), tt.expected)
			}
		})
	}
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "debug level",
			input:    `{"level":"debug","message":"test"}`,
			expected: "debug",
		},
		{
			name:     "info level",
			input:    `{"level":"info","message":"test"}`,
			expected: "info",
		},
		{
			name:     "warn level",
			input:    `{"level":"warn","message":"test"}`,
			expected: "warn",
		},
		{
			name:     "error level",
			input:    `{"level":"error","message":"test"}`,
			expected: "error",
		},
		{
			name:     "fatal level",
			input:    `{"level":"fatal","message":"test"}`,
			expected: "fatal",
		},
		{
			name:     "unknown level",
			input:    `{"level":"unknown","message":"test"}`,
			expected: "",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "malformed json",
			input:    `{"level":"debug"`,
			expected: "debug",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLogLevel([]byte(tt.input))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestColorWriter_Colorize(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		message  []byte
		expected string
	}{
		{
			name:     "debug level",
			level:    "debug",
			message:  []byte("test"),
			expected: "test",
		},
		{
			name:     "unknown level",
			level:    "unknown",
			message:  []byte("test"),
			expected: "test",
		},
		{
			name:     "empty message",
			level:    "debug",
			message:  []byte(""),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw := &ColorWriter{
				colors: map[string]string{
					"debug": "\033[36m",
					"reset": "\033[0m",
				},
			}

			result := cw.colorize(tt.level, tt.message)
			if tt.level == "debug" {
				assert.Contains(t, string(result), "\033[36m")
				assert.Contains(t, string(result), "\033[0m")
			} else {
				assert.Equal(t, tt.message, result)
			}
		})
	}
}

func TestSetWindowsColors(t *testing.T) {
	cw := &ColorWriter{}
	cw.setWindowsColors()

	if cw.colors["debug"] != "\033[36m" {
		t.Errorf("expected cyan for debug, got %s", cw.colors["debug"])
	}
}
