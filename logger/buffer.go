// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package logger

import (
	"bytes"
	"io"
	"sync"
)

// NewBufferedWriter creates a new buffered writer that supports pause/resume operations.
// This is useful for temporarily buffering log output and flushing it later.
func NewBufferedWriter(out io.Writer) io.Writer {
	return &BufferedWriter{
		out: out,
	}
}

// BufferedWriter implements a thread-safe buffered writer with pause/resume capabilities.
// It can be used to temporarily buffer log output and flush it later, which is useful
// for operations that need to collect logs before displaying them.
type BufferedWriter struct {
	out    io.Writer
	buf    bytes.Buffer
	paused bool
	lock   sync.RWMutex
}

// Pause stops writing to the output and starts buffering the data.
// All writes after calling Pause will be stored in the buffer until Resume is called.
func (bw *BufferedWriter) Pause() {
	bw.lock.Lock()
	defer bw.lock.Unlock()

	bw.paused = true
}

// Resume flushes the buffered data to the output and resumes normal operation.
// If the writer is not paused, this is a no-op.
func (bw *BufferedWriter) Resume() {
	bw.lock.Lock()
	defer bw.lock.Unlock()

	if !bw.paused {
		return
	}

	bw.paused = false
	if bw.buf.Len() > 0 {
		_, _ = bw.out.Write(bw.buf.Bytes())
		bw.buf.Reset()
	}
}

// Write implements the io.Writer interface.
// If the writer is paused, it writes to the buffer.
// Otherwise, it writes directly to the output.
func (bw *BufferedWriter) Write(p []byte) (n int, err error) {
	bw.lock.RLock()
	defer bw.lock.RUnlock()

	if bw.paused {
		return bw.buf.Write(p)
	}

	return bw.out.Write(p)
}
