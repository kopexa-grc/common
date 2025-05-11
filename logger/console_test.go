// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsoleFormatFunctions(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		expected string
	}{
		{
			name:     "TRACE level",
			level:    TRACE,
			expected: "TRC",
		},
		{
			name:     "DEBUG level",
			level:    DEBUG,
			expected: "DBG",
		},
		{
			name:     "INFO level",
			level:    INFO,
			expected: "-",
		},
		{
			name:     "WARN level",
			level:    WARN,
			expected: "WRN",
		},
		{
			name:     "ERROR level",
			level:    ERROR,
			expected: "ERR",
		},
		{
			name:     "FATAL level",
			level:    FATAL,
			expected: "FTL",
		},
		{
			name:     "PANIC level",
			level:    PANIC,
			expected: "PNC",
		},
		{
			name:     "Unknown level",
			level:    "UNKNOWN",
			expected: "UNK",
		},
		{
			name:     "Nil level",
			level:    "",
			expected: "UNK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test consoleFormatLevelNoColor
			formatter := consoleFormatLevelNoColor()
			result := formatter(tt.level)
			assert.Equal(t, tt.expected, result)

			// Test consoleFormatLevel
			formatter = consoleFormatLevel()
			result = formatter(tt.level)
			// We can't test the exact color output, but we can verify it's not empty
			assert.NotEmpty(t, result)
		})
	}
}

func TestConsoleDefaultFormatFunctions(t *testing.T) {
	t.Run("FormatCaller", func(t *testing.T) {
		formatter := consoleDefaultFormatCaller()
		result := formatter("test/caller.go:123")
		assert.Contains(t, result, "test/caller.go:123")
	})

	t.Run("FormatMessage", func(t *testing.T) {
		result := consoleDefaultFormatMessage("test message")
		assert.Equal(t, "test message", result)

		result = consoleDefaultFormatMessage(nil)
		assert.Equal(t, "", result)
	})

	t.Run("FormatFieldName", func(t *testing.T) {
		formatter := consoleDefaultFormatFieldName()
		result := formatter("field")
		assert.Contains(t, result, "field=")
	})

	t.Run("FormatFieldValue", func(t *testing.T) {
		result := consoleDefaultFormatFieldValue("value")
		assert.Equal(t, "value", result)
	})

	t.Run("FormatErrFieldName", func(t *testing.T) {
		formatter := consoleDefaultFormatErrFieldName()
		result := formatter("error")
		assert.Contains(t, result, "error=")
	})

	t.Run("FormatErrFieldValue", func(t *testing.T) {
		formatter := consoleDefaultFormatErrFieldValue()
		result := formatter("error value")
		assert.Contains(t, result, "error value")
	})
}

func TestNewConsoleWriter(t *testing.T) {
	tests := []struct {
		name     string
		compact  bool
		level    string
		message  string
		expected string
	}{
		{
			name:     "Compact mode",
			compact:  true,
			level:    INFO,
			message:  "test message",
			expected: "- test message",
		},
		{
			name:     "Non-compact mode",
			compact:  false,
			level:    INFO,
			message:  "test message",
			expected: "test message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewConsoleWriter(&buf, tt.compact)

			switch tt.level {
			case TRACE:
				logger.Trace().Msg(tt.message)
			case DEBUG:
				logger.Debug().Msg(tt.message)
			case INFO:
				logger.Info().Msg(tt.message)
			case WARN:
				logger.Warn().Msg(tt.message)
			case ERROR:
				logger.Error().Msg(tt.message)
			case FATAL:
				logger.Fatal().Msg(tt.message)
			case PANIC:
				logger.Panic().Msg(tt.message)
			}

			output := buf.String()
			assert.Contains(t, output, tt.message)
		})
	}
}
