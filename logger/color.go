// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// ColorWriter implements a writer that adds color to the output.
// It supports different color schemes for different operating systems
// and can be disabled for non-terminal outputs.
type ColorWriter struct {
	out    io.Writer
	colors map[string]string
}

// NewColorWriter creates a new color writer with the appropriate color scheme
// for the current operating system.
func NewColorWriter(out io.Writer) *ColorWriter {
	if out == nil {
		out = os.Stdout
	}

	cw := &ColorWriter{
		out:    out,
		colors: make(map[string]string),
	}

	// Set colors based on OS
	if runtime.GOOS == "windows" {
		cw.setWindowsColors()
	} else {
		cw.setUnixColors()
	}

	return cw
}

// colorize wraps the message with the color code for the given level, if available.
func (cw *ColorWriter) colorize(level string, msg []byte) []byte {
	if color, ok := cw.colors[level]; ok {
		return []byte(fmt.Sprintf("%s%s%s", color, msg, cw.colors["reset"]))
	}

	return msg
}

// Write implements the io.Writer interface.
// It adds color to the output based on the log level.
func (cw *ColorWriter) Write(p []byte) (n int, err error) {
	// Check if output is a terminal
	if !isTerminal(cw.out) {
		return cw.out.Write(p)
	}

	// Add color based on log level
	level := getLogLevel(p)

	return cw.out.Write(cw.colorize(level, p))
}

// setWindowsColors sets the color scheme for Windows.
func (cw *ColorWriter) setWindowsColors() {
	cw.colors = map[string]string{
		"debug": "\033[36m", // Cyan
		"info":  "\033[32m", // Green
		"warn":  "\033[33m", // Yellow
		"error": "\033[31m", // Red
		"fatal": "\033[35m", // Magenta
		"reset": "\033[0m",  // Reset
	}
}

// setUnixColors sets the color scheme for Unix-like systems.
func (cw *ColorWriter) setUnixColors() {
	cw.colors = map[string]string{
		"debug": "\033[36m", // Cyan
		"info":  "\033[32m", // Green
		"warn":  "\033[33m", // Yellow
		"error": "\033[31m", // Red
		"fatal": "\033[35m", // Magenta
		"reset": "\033[0m",  // Reset
	}
}

// getLogLevel extracts the log level from the log message.
func getLogLevel(p []byte) string {
	msg := string(p)
	if strings.Contains(msg, "\"level\":\"debug\"") {
		return "debug"
	}

	if strings.Contains(msg, "\"level\":\"info\"") {
		return "info"
	}

	if strings.Contains(msg, "\"level\":\"warn\"") {
		return "warn"
	}

	if strings.Contains(msg, "\"level\":\"error\"") {
		return "error"
	}

	if strings.Contains(msg, "\"level\":\"fatal\"") {
		return "fatal"
	}

	return ""
}

// isTerminal checks if the output is a terminal.
func isTerminal(w io.Writer) bool {
	if f, ok := w.(*os.File); ok {
		return isTerminalFile(f)
	}

	return false
}

// isTerminalFile checks if the file is a terminal.
func isTerminalFile(f *os.File) bool {
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}
