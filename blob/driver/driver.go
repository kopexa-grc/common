// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package driver

import (
	"context"
	"io"
	"time"
)

// CopyOptions is the options for the Copy method.
type CopyOptions struct {
	// BeforeCopy is a callback that will be called before each call to the
	//
	BeforeCopy func(asFunc func(any) bool) error
}

// Bucket provides read, write and delete operations on objects within it on the
// blob service.
type Bucket interface {
	// Delete deletes the object associated with key. If the specified object does
	// not exist, Delete must return an error for which ErrorCode returns
	// gcerrors.NotFound.
	Delete(ctx context.Context, key string) error
	// SignedURL returns a URL that can be used to GET the blob for the duration
	// specified in opts.Expiry. opts is guaranteed to be non-nil.
	// If not supported, return an error for which ErrorCode returns
	// kerrs.Unimplemented.
	SignedURL(ctx context.Context, key string, opts *SignedURLOptions) (string, error)

	// Copy copies the object from srcKey to dstKey.
	//
	// If the source object does not exist, Copy must return an error for which
	// ErrorCode returns kerr.NotFound.
	//
	// If the destination object already exists, it should be overwritten.
	//
	// opts is guaranteed to be non-nil.
	Copy(ctx context.Context, srcKey, dstKey string, opts *CopyOptions) error

	// NewRangeReader returns a Reader that reads part of an object, reading at
	// most length bytes starting at the given offset. If length is negative, it
	// will read until the end of the object. If the specified object does not
	// exist, NewRangeReader must return an error for which ErrorCode returns
	// gcerrors.NotFound.
	// opts is guaranteed to be non-nil.
	//
	// The returned Reader *may* also implement Downloader if the underlying
	// implementation can take advantage of that. The Download call is guaranteed
	// to be the only call to the Reader. For such readers, offset will always
	// be 0 and length will always be -1.
	NewRangeReader(ctx context.Context, key string, offset, length int64, opts *ReaderOptions) (Reader, error)

	// NewTypedWriter returns Writer that writes to an object associated with key.
	//
	// A new object will be created unless an object with this key already exists.
	// Otherwise any previous object with the same key will be replaced.
	// The object may not be available (and any previous object will remain)
	// until Close has been called.
	//
	// contentType sets the MIME type of the object to be written.
	// opts is guaranteed to be non-nil.
	//
	// The caller must call Close on the returned Writer when done writing.
	//
	// Implementations should abort an ongoing write if ctx is later canceled,
	// and do any necessary cleanup in Close. Close should then return ctx.Err().
	//
	// The returned Writer *may* also implement Uploader if the underlying
	// implementation can take advantage of that. The Upload call is guaranteed
	// to be the only non-Close call to the Writer..
	NewTypedWriter(ctx context.Context, key, contentType string, opts *WriterOptions) (Writer, error)
}

// SignedURLOptions sets options for SignedURL.
type SignedURLOptions struct {
	// Expiry sets how long the returned URL is valid for. It is guaranteed to be > 0.
	Expiry time.Duration
	// Method is the HTTP method that can be used on the URL; one of "GET", "PUT",
	// or "DELETE". Drivers must implement all 3.
	Method string

	// ContentType specifies the Content-Type HTTP header the user agent is
	// permitted to use in the PUT request. It must match exactly. See
	// EnforceAbsentContentType for behavior when ContentType is the empty string.
	// If this field is not empty and the bucket cannot enforce the Content-Type
	// header, it must return an Unimplemented error.
	//
	// This field will not be set for any non-PUT requests.
	ContentType string

	// BeforeSign is a callback that will be called before each call to the
	// the underlying service's sign functionality.
	// asFunc converts its argument to driver-specific types.
	BeforeSign func(asFunc func(any) bool) error
}

// ReaderOptions controls Reader behaviors.
type ReaderOptions struct {
	// BeforeRead is a callback that must be called exactly once before
	// any data is read, unless NewRangeReader returns an error before then, in
	// which case it should not be called at all.
	// asFunc allows drivers to expose driver-specific types;
	// see Bucket.As for more details.
	BeforeRead func(asFunc func(any) bool) error
}

// Reader reads an object from the blob.
type Reader interface {
	io.ReadCloser

	// Attributes returns a subset of attributes about the blob.
	// The portable type will not modify the returned ReaderAttributes.
	Attributes() *ReaderAttributes

	// As allows drivers to expose driver-specific types;
	// see Bucket.As for more details.
	As(any) bool
}

// ReaderAttributes contains a subset of attributes about a blob that are
// accessible from Reader.
type ReaderAttributes struct {
	// ContentType is the MIME type of the blob object. It must not be empty.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	ContentType string
	// ModTime is the time the blob object was last modified.
	ModTime time.Time
	// Size is the size of the object in bytes.
	Size int64
}

// Downloader has an optional extra method for readers.
// It is similar to io.WriteTo, but without the count of bytes returned.
type Downloader interface {
	// Download is similar to io.WriteTo, but without the count of bytes returned.
	Download(w io.Writer) error
}

// Uploader has an optional extra method for writers.
type Uploader interface {
	// Upload is similar to io.ReadFrom, but without the count of bytes returned.
	Upload(r io.Reader) error
}

// Writer writes an object to the blob.
type Writer interface {
	io.WriteCloser
}

// WriterOptions controls behaviors of Writer.
type WriterOptions struct {
	// BufferSize changes the default size in byte of the maximum part Writer can
	// write in a single request, if supported. Larger objects will be split into
	// multiple requests.
	BufferSize int
	// MaxConcurrency changes the default concurrency for uploading parts.
	MaxConcurrency int
	// CacheControl specifies caching attributes that services may use
	// when serving the blob.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
	CacheControl string
	// ContentDisposition specifies whether the blob content is expected to be
	// displayed inline or as an attachment.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
	ContentDisposition string
	// ContentEncoding specifies the encoding used for the blob's content, if any.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Encoding
	ContentEncoding string
	// ContentLanguage specifies the language used in the blob's content, if any.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Language
	ContentLanguage string
	// ContentMD5 is used as a message integrity check.
	// The portable type checks that the MD5 hash of the bytes written matches
	// ContentMD5.
	// If len(ContentMD5) > 0, driver implementations may pass it to their
	// underlying network service to guarantee the integrity of the bytes in
	// transit.
	ContentMD5 []byte
	// Metadata holds key/value strings to be associated with the blob.
	// Keys are guaranteed to be non-empty and lowercased.
	Metadata map[string]string
	// When true, the driver should attempt to disable any automatic
	// content-type detection that the provider applies on writes with an
	// empty ContentType.
	DisableContentTypeDetection bool
	// BeforeWrite is a callback that must be called exactly once before
	// any data is written, unless NewTypedWriter returns an error, in
	// which case it should not be called.
	// asFunc allows drivers to expose driver-specific types;
	// see Bucket.As for more details.
	BeforeWrite func(asFunc func(any) bool) error

	// IfNotExist is used for conditional writes.
	// When set to true, if a blob exists for the same key in the bucket, the write operation
	// won't take place.
	IfNotExist bool
}
