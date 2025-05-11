// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package logger

import (
	"bytes"
	"io"
	"testing"
)

func TestBufferedWriter(t *testing.T) {
	tests := []struct {
		name     string
		write    func(w io.Writer)
		expected string
	}{
		{
			name: "direct write",
			write: func(w io.Writer) {
				_, _ = w.Write([]byte("test"))
			},
			expected: "test",
		},
		{
			name: "paused write",
			write: func(w io.Writer) {
				bw := w.(*BufferedWriter)
				bw.Pause()
				_, _ = w.Write([]byte("test"))
				bw.Resume()
			},
			expected: "test",
		},
		{
			name: "multiple writes",
			write: func(w io.Writer) {
				bw := w.(*BufferedWriter)
				bw.Pause()
				_, _ = w.Write([]byte("test1"))
				_, _ = w.Write([]byte("test2"))
				bw.Resume()
			},
			expected: "test1test2",
		},
		{
			name: "nested pause/resume",
			write: func(w io.Writer) {
				bw := w.(*BufferedWriter)
				bw.Pause()
				_, _ = w.Write([]byte("test1"))
				bw.Pause()
				_, _ = w.Write([]byte("test2"))
				bw.Resume()
				_, _ = w.Write([]byte("test3"))
				bw.Resume()
			},
			expected: "test1test2test3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bw := NewBufferedWriter(&buf)
			tt.write(bw)

			if got := buf.String(); got != tt.expected {
				t.Errorf("BufferedWriter.Write() = %v, want %v", got, tt.expected)
			}
		})
	}
}
