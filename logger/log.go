// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package logger provides a flexible and enterprise-ready logging solution.
// It supports multiple output formats (JSON, console), environments (local, Docker, Kubernetes),
// and includes features like log buffering, color support, and structured logging.
//
// Key features:
// - Multiple output formats (JSON, console, GCP-compatible JSON)
// - Environment-aware configuration (local, Docker, Kubernetes)
// - Buffered logging with pause/resume capabilities
// - Color support for console output
// - Structured logging with field support
// - Configurable log levels
// - Thread-safe operations
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogOutputWriter is the default output writer for logs.
// It uses a buffered writer to support pause/resume operations.
var LogOutputWriter = NewBufferedWriter(os.Stderr)

// Debug indicates if the application is running in debug mode.
// This can be used to enable additional debug logging.
var Debug bool

// SetWriter configures a custom writer for the global logger.
// This is useful for redirecting logs to different outputs or implementing custom writers.
func SetWriter(w io.Writer) {
	log.Logger = log.Output(w)
}

// UseJSONLogging configures the global logger to output JSON format.
// This is suitable for machine processing and log aggregation systems.
func UseJSONLogging(out io.Writer) {
	log.Logger = zerolog.New(out).With().Timestamp().Logger()
}

// UseGCPJSONLogging configures the global logger to output GCP-compatible JSON format.
// This includes specific field names and formats that Google Cloud Platform expects.
func UseGCPJSONLogging(out io.Writer) {
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log.Logger = zerolog.New(out).With().Timestamp().Logger()
}

// CliLogger configures the global logger for console output with colors.
// This is suitable for local development and debugging.
func CliLogger() {
	log.Logger = NewConsoleWriter(LogOutputWriter, false)
}

// CliCompactLogger configures the global logger for compact console output.
// This is useful when space is limited or for dense log displays.
func CliCompactLogger(out io.Writer) {
	log.Logger = NewConsoleWriter(out, true)
}

// StandardZerologLogger configures the global logger to use zerolog's standard console writer.
// This provides a basic console output format with timestamps.
func StandardZerologLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

// Set configures the global log level.
// Supported levels: "error", "warn", "info", "debug", "trace".
// If an invalid level is provided, it defaults to "info" and logs an error.
func Set(level string) {
	switch level {
	case ERROR:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case WARN:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case INFO:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case DEBUG:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case TRACE:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		log.Error().Str("level", level).Msg("unknown log level, defaulting to info")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}

// GetLevel returns the current global log level as a string.
func GetLevel() string {
	return zerolog.GlobalLevel().String()
}

// InitTestEnv configures the logger for test environments.
// It sets the log level to debug and disables colors for consistent test output.
func InitTestEnv() {
	Set(DEBUG)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
}

// GetEnvLogLevel determines the log level from environment variables.
// It checks for DEBUG and TRACE environment variables.
// Returns the level and true if a level was found, empty string and false otherwise.
func GetEnvLogLevel() (string, bool) {
	level := ""
	ok := false

	if os.Getenv("DEBUG") == "true" || os.Getenv("DEBUG") == "1" {
		level = DEBUG
		ok = true
	}

	if os.Getenv("TRACE") == "true" || os.Getenv("TRACE") == "1" {
		level = TRACE
		ok = true
	}

	return level, ok
}
