// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package logger

import (
	"bytes"
	"os"
	"testing"
)

func TestSetWriter(t *testing.T) {
	var buf bytes.Buffer

	SetWriter(&buf)

	log := LogOutputWriter

	_, err := log.Write([]byte("test\n"))
	if err != nil {
		t.Errorf("SetWriter: unexpected error: %v", err)
	}
}

func TestUseJSONLogging(t *testing.T) {
	var buf bytes.Buffer

	UseJSONLogging(&buf)

	log := LogOutputWriter

	_, err := log.Write([]byte("test\n"))
	if err != nil {
		t.Errorf("UseJSONLogging: unexpected error: %v", err)
	}
}

func TestUseGCPJSONLogging(t *testing.T) {
	var buf bytes.Buffer

	UseGCPJSONLogging(&buf)

	log := LogOutputWriter

	_, err := log.Write([]byte("test\n"))
	if err != nil {
		t.Errorf("UseGCPJSONLogging: unexpected error: %v", err)
	}
}

func TestCliCompactLogger(t *testing.T) {
	var buf bytes.Buffer

	CliCompactLogger(&buf)

	log := LogOutputWriter

	_, err := log.Write([]byte("test\n"))
	if err != nil {
		t.Errorf("CliCompactLogger: unexpected error: %v", err)
	}
}

func TestStandardZerologLogger(t *testing.T) {
	StandardZerologLogger()

	log := LogOutputWriter

	_, err := log.Write([]byte("test\n"))
	if err != nil {
		t.Errorf("StandardZerologLogger: unexpected error: %v", err)
	}
}

func TestSetLevels(_ *testing.T) {
	levels := []string{"error", "warn", "info", "debug", "trace", "", "invalid"}
	for _, level := range levels {
		Set(level)
	}
}

func TestGetLevel(t *testing.T) {
	Set("debug")

	if GetLevel() != "debug" {
		t.Errorf("expected debug, got %s", GetLevel())
	}
}

func TestInitTestEnv(t *testing.T) {
	InitTestEnv()

	if GetLevel() != "debug" {
		t.Errorf("expected debug, got %s", GetLevel())
	}
}

func TestGetEnvLogLevel(t *testing.T) {
	os.Setenv("DEBUG", "true")

	level, ok := GetEnvLogLevel()
	if !ok || level != "debug" {
		t.Errorf("expected debug, got %s", level)
	}

	os.Unsetenv("DEBUG")

	os.Setenv("TRACE", "true")

	level, ok = GetEnvLogLevel()
	if !ok || level != "trace" {
		t.Errorf("expected trace, got %s", level)
	}

	os.Unsetenv("TRACE")

	level, ok = GetEnvLogLevel()
	if ok || level != "" {
		t.Errorf("expected empty, got %s", level)
	}
}
