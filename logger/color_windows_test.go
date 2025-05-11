// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

//go:build windows
// +build windows

package logger

import (
	"os"
	"testing"
)

func TestEnableVirtualTerminal(t *testing.T) {
	// Save original stdout
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Set stdout to the temporary file
	os.Stdout = tmpFile

	// Test enableVirtualTerminal
	err = enableVirtualTerminal()
	if err != nil {
		t.Errorf("enableVirtualTerminal() error = %v", err)
	}

	// Verify that the file is writable
	_, err = tmpFile.WriteString("test")
	if err != nil {
		t.Errorf("Failed to write to temp file: %v", err)
	}
}
