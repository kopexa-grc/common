// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

//go:build windows
// +build windows

package logger

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleMode             = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode             = kernel32.NewProc("SetConsoleMode")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

const (
	// Windows console mode flags
	enableVirtualTerminalProcessing = 0x0004
)

// enableVirtualTerminal enables virtual terminal processing for Windows console.
// This is required for ANSI escape sequences to work properly.
func enableVirtualTerminal() error {
	var mode uint32
	handle := syscall.Handle(os.Stdout.Fd())

	// Get current console mode
	ret, _, err := procGetConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return err
	}

	// Enable virtual terminal processing
	mode |= enableVirtualTerminalProcessing
	ret, _, err = procSetConsoleMode.Call(uintptr(handle), uintptr(mode))
	if ret == 0 {
		return err
	}

	return nil
}

// init initializes the Windows console for color output.
func init() {
	_ = enableVirtualTerminal()
}
