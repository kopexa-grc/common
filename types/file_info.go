// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package types provides core data structures used throughout the application.
package types

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	// MaxFileSize is the maximum allowed file size in bytes (100MB)
	MaxFileSize = 100 * 1024 * 1024
	KB          = 1024
	MB          = KB * 1024
	GB          = MB * 1024

	// AllowedMimeTypes contains the list of allowed MIME types
	AllowedMimeTypes = "application/pdf,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,image/jpeg,image/png,text/plain"
)

var (
	ErrFileSizeExceeded      = errors.New("file size exceeds maximum allowed size")
	ErrContentTypeNotAllowed = errors.New("content type is not allowed")
	ErrContentTypeMismatch   = errors.New("content type does not match file extension")
)

// FileInfo represents metadata for a file upload in a document.
// It contains essential information about the file including its location,
// name, size, and content type.
//
// @swagger:model FileInfo
type FileInfo struct {
	// Path is the internal storage path of the file.
	Path string `json:"path"`

	// Name is the original filename of the uploaded file.
	// @example "policy.pdf"
	Name string `json:"name" validate:"required" example:"policy.pdf"`

	// URL is the publicly accessible URL of the file.
	// @example "https://example.com/files/policy.pdf"
	URL string `json:"url" validate:"required,url" example:"https://example.com/files/policy.pdf"`

	// Size is the size of the file in bytes.
	// @example 1048576
	Size int64 `json:"size" validate:"required" example:"1048576"`

	// ContentType is the MIME type of the file.
	// @example "application/pdf"
	ContentType string `json:"contentType" validate:"required" example:"application/pdf"`
}

// Validate checks if the FileInfo is valid.
// It verifies the file size and content type.
//
// Returns:
//   - error: If the FileInfo is invalid
func (d FileInfo) Validate() error {
	if d.Size > MaxFileSize {
		return fmt.Errorf("%w: %d bytes", ErrFileSizeExceeded, MaxFileSize)
	}

	allowedTypes := strings.Split(AllowedMimeTypes, ",")
	if !slices.Contains(allowedTypes, d.ContentType) {
		return fmt.Errorf("%w: %s", ErrContentTypeNotAllowed, d.ContentType)
	}

	// Validate that the content type matches the file extension
	ext := strings.ToLower(d.Name[strings.LastIndex(d.Name, ".")+1:])
	mimeType := mime.TypeByExtension("." + ext)

	if mimeType != "" && mimeType != d.ContentType {
		return fmt.Errorf("%w: %s vs %s", ErrContentTypeMismatch, d.ContentType, ext)
	}

	return nil
}

// String returns a human-readable string representation of the FileInfo.
// It formats the file information in a clear, concise way.
//
// Example:
//
//	File: policy.pdf (1.0 MB)
//	Type: application/pdf
//	URL: https://example.com/files/policy.pdf
func (d FileInfo) String() string {
	if d == (FileInfo{}) {
		return "<empty file info>"
	}

	var sizeStr string

	switch {
	case d.Size == 0:
		sizeStr = "0 B"
	case d.Size < KB:
		sizeStr = fmt.Sprintf("%d B", d.Size)
	case d.Size < MB:
		sizeStr = fmt.Sprintf("%.1f KB", float64(d.Size)/KB)
	case d.Size < GB:
		sizeStr = fmt.Sprintf("%.1f MB", float64(d.Size)/MB)
	default:
		sizeStr = fmt.Sprintf("%.1f GB", float64(d.Size)/GB)
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("File: %s (%s)\n", d.Name, sizeStr))
	sb.WriteString(fmt.Sprintf("Type: %s\n", d.ContentType))
	sb.WriteString(fmt.Sprintf("URL: %s", d.URL))

	return sb.String()
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for FileInfo.
// It allows FileInfo to be used as a GraphQL scalar type.
//
// Parameters:
//   - v: The value to unmarshal
//
// Returns:
//   - error: If unmarshaling fails
func (d *FileInfo) UnmarshalGQL(v any) error {
	return unmarshalGQLJSON(v, d)
}

// MarshalGQL implements the graphql.Marshaler interface for FileInfo.
// It allows FileInfo to be used as a GraphQL scalar type.
//
// Parameters:
//   - w: The writer to write the FileInfo to
func (d FileInfo) MarshalGQL(w io.Writer) {
	if err := marshalGQLJSON(w, d); err != nil {
		log.Error().Err(err).Msg("failed to marshal file info to GraphQL")
	}
}
