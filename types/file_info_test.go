// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileInfo_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   FileInfo
		wantErr bool
	}{
		{
			name: "valid file info",
			input: FileInfo{
				Name:        "test.pdf",
				Path:        "/path/to/test.pdf",
				URL:         "https://example.com/test.pdf",
				Size:        1024,
				ContentType: "application/pdf",
			},
		},
		{
			name: "file too large",
			input: FileInfo{
				Name:        "test.pdf",
				Path:        "/path/to/test.pdf",
				URL:         "https://example.com/test.pdf",
				Size:        MaxFileSize + 1,
				ContentType: "application/pdf",
			},
			wantErr: true,
		},
		{
			name: "invalid content type",
			input: FileInfo{
				Name:        "test.exe",
				Path:        "/path/to/test.exe",
				URL:         "https://example.com/test.exe",
				Size:        1024,
				ContentType: "application/x-executable",
			},
			wantErr: true,
		},
		{
			name: "content type mismatch",
			input: FileInfo{
				Name:        "test.pdf",
				Path:        "/path/to/test.pdf",
				URL:         "https://example.com/test.pdf",
				Size:        1024,
				ContentType: "image/jpeg",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestFileInfo_String(t *testing.T) {
	tests := []struct {
		name  string
		input FileInfo
		want  string
	}{
		{
			name:  "empty file info",
			input: FileInfo{},
			want:  "<empty file info>",
		},
		{
			name: "valid file info",
			input: FileInfo{
				Name:        "test.pdf",
				Path:        "/path/to/test.pdf",
				URL:         "https://example.com/test.pdf",
				Size:        1024,
				ContentType: "application/pdf",
			},
			want: "File: test.pdf (1.0 KB)\nType: application/pdf\nURL: https://example.com/test.pdf",
		},
		{
			name: "large file",
			input: FileInfo{
				Name:        "large.pdf",
				Path:        "/path/to/large.pdf",
				URL:         "https://example.com/large.pdf",
				Size:        1024 * 1024 * 2, // 2MB
				ContentType: "application/pdf",
			},
			want: "File: large.pdf (2.0 MB)\nType: application/pdf\nURL: https://example.com/large.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.String())
		})
	}
}

func TestFileInfo_MarshalGQL(t *testing.T) {
	tests := []struct {
		name  string
		input FileInfo
		want  string
	}{
		{
			name:  "empty file info",
			input: FileInfo{},
			want:  `{"path":"","name":"","url":"","size":0,"contentType":""}`,
		},
		{
			name: "valid file info",
			input: FileInfo{
				Name:        "test.pdf",
				Path:        "/path/to/test.pdf",
				URL:         "https://example.com/test.pdf",
				Size:        1024,
				ContentType: "application/pdf",
			},
			want: `{"path":"/path/to/test.pdf","name":"test.pdf","url":"https://example.com/test.pdf","size":1024,"contentType":"application/pdf"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			tt.input.MarshalGQL(&buf)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestFileInfo_UnmarshalGQL(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    FileInfo
		wantErr bool
	}{
		{
			name: "valid file info",
			input: map[string]interface{}{
				"name":        "test.pdf",
				"path":        "/path/to/test.pdf",
				"url":         "https://example.com/test.pdf",
				"size":        1024,
				"contentType": "application/pdf",
			},
			want: FileInfo{
				Name:        "test.pdf",
				Path:        "/path/to/test.pdf",
				URL:         "https://example.com/test.pdf",
				Size:        1024,
				ContentType: "application/pdf",
			},
		},
		{
			name:    "invalid type",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fi FileInfo
			err := fi.UnmarshalGQL(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, fi)
		})
	}
}
