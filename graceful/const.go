// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package graceful

import "time"

// Log message constants
const (
	LogMsgRecoveredFromPanic = "Recovered from panic"
	LogMsgTriggeringShutdown = "Triggering shutdown from signal"
	LogMsgShuttingDown       = "Shutting down..."
	LogMsgShutdownTimedOut   = "Shutdown timed out"
	LogMsgShutdownFailed     = "Shutdown failed"
	LogMsgShutdownFinished   = "Shutdown finished"
)

// Log field constants
const (
	LogFieldComponent = "component"
	LogFieldPanic     = "panic"
	LogFieldSignal    = "signal"
	LogFieldTarget    = "target"
	LogFieldTimeout   = "timeout"
)

// Timeout constants
const (
	MinShutdownTimeout = 1 * time.Second
	MaxShutdownTimeout = 30 * time.Second
)

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
)
