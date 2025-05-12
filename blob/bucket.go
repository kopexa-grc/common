// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob

import (
	"context"
	"net/http"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/kopexa-grc/common/blob/driver"
	kerr "github.com/kopexa-grc/common/errors"
)

// Bucket provides an easy and portable way to interact with blobs
// within a "bucket", including read, write, and list operations.
// To create a Bucket, use constructors found in driver subpackages.
type Bucket struct {
	b driver.Bucket

	// mu protects the closed variable.
	// Read locks are kept to allow holding a read lock for long-running calls,
	// and thereby prevent closing until a call finishes.
	mu     sync.RWMutex
	closed bool
}

// Delete deletes the blob stored at key.
//
// If the blob does not exist, Delete returns an error for which
// gcerrors.Code will return gcerrors.NotFound.
func (b *Bucket) Delete(ctx context.Context, key string) (err error) {
	if !utf8.ValidString(key) {
		return kerr.Newf(kerr.InvalidArgument, nil, "blob: Delete key must be a valid UTF-8 string: %q", key)
	}

	if key == "" {
		return kerr.Newf(kerr.InvalidArgument, nil, "blob: Delete key must be a non-empty string")
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return errClosed
	}

	return b.b.Delete(ctx, key)
}

// SignedURLOptions sets options for SignedURL.
type SignedURLOptions struct {
	// Expiry sets how long the returned URL is valid for.
	// Defaults to DefaultSignedURLExpiry.
	Expiry time.Duration

	// Method is the HTTP method that can be used on the URL; one of "GET", "PUT",
	// or "DELETE". Defaults to "GET".
	Method string
	// ContentType specifies the Content-Type HTTP header the user agent is
	// permitted to use in the PUT request. It must match exactly. See
	// EnforceAbsentContentType for behavior when ContentType is the empty string.
	// If a bucket does not implement this verification, then it returns an
	// Unimplemented error.
	//
	// Must be empty for non-PUT requests.
	ContentType string

	// BeforeSign is a callback that will be called before each call to the
	// the underlying service's sign functionality.
	// asFunc converts its argument to driver-specific types.
	// See https://gocloud.dev/concepts/as/ for background information.
	BeforeSign func(asFunc func(any) bool) error
}

// SignedURL returns a URL that can be used to GET (default), PUT or DELETE
// the blob for the duration specified in opts.Expiry.
//
// A nil SignedURLOptions is treated the same as the zero value.
//
// It is valid to call SignedURL for a key that does not exist.
//
// If the driver does not support this functionality, SignedURL
// will return an error for which gcerrors.Code will return gcerrors.Unimplemented.
func (b *Bucket) SignedURL(ctx context.Context, key string, opts *SignedURLOptions) (string, error) {
	if !utf8.ValidString(key) {
		return "", kerr.Newf(kerr.InvalidArgument, nil, "blob: SignedURL key must be a valid UTF-8 string: %q", key)
	}
	dopts := new(driver.SignedURLOptions)
	if opts == nil {
		opts = new(SignedURLOptions)
	}

	switch {
	case opts.Expiry < 0:
		return "", kerr.Newf(kerr.InvalidArgument, nil, "blob: SignedURL expiry must be non-negative: %q", opts.Expiry)
	case opts.Expiry == 0:
		dopts.Expiry = DefaultSignedURLExpiry
	default:
		dopts.Expiry = opts.Expiry
	}

	switch opts.Method {
	case "":
		dopts.Method = http.MethodGet
	case http.MethodGet, http.MethodPut, http.MethodDelete:
		dopts.Method = opts.Method
	default:
		return "", kerr.Newf(kerr.InvalidArgument, nil, "blob: SignedURL method must be one of GET, PUT, or DELETE: %q", opts.Method)
	}

	if opts.ContentType != "" && opts.Method != http.MethodPut {
		return "", kerr.Newf(kerr.InvalidArgument, nil, "blob: SignedURL ContentType must be empty for non-PUT requests: %q", opts.ContentType)
	}

	dopts.ContentType = opts.ContentType
	dopts.BeforeSign = opts.BeforeSign

	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return "", errClosed
	}

	url, err := b.b.SignedURL(ctx, key, dopts)
	if err != nil {
		return "", err
	}

	return url, nil
}

// CopyOptions sets options for Copy.
type CopyOptions struct {
	// BeforeCopy is a callback that will be called before the copy is
	// initiated.
	//
	// asFunc converts its argument to driver-specific types.
	BeforeCopy func(asFunc func(any) bool) error
}

// Copy the blob stored at srcKey to dstKey.
// A nil CopyOptions is treated the same as the zero value.
//
// If the source blob does not exist, Copy returns an error for which
// kerr.Code will return kerr.NotFound.
//
// If the destination blob already exists, it is overwritten.
func (b *Bucket) Copy(ctx context.Context, dstKey, srcKey string, opts *CopyOptions) (err error) {
	if !utf8.ValidString(dstKey) {
		return kerr.Newf(kerr.InvalidArgument, nil, "blob: Copy dstKey must be a valid UTF-8 string: %q", dstKey)
	}
	if !utf8.ValidString(srcKey) {
		return kerr.Newf(kerr.InvalidArgument, nil, "blob: Copy srcKey must be a valid UTF-8 string: %q", srcKey)
	}

	if opts == nil {
		opts = new(CopyOptions)
	}

	dopts := &driver.CopyOptions{
		BeforeCopy: opts.BeforeCopy,
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return errClosed
	}

	return b.b.Copy(ctx, dstKey, srcKey, dopts)
}

// NewBucketForTest creates a Bucket with a mock driver for testing purposes.
func NewBucketForTest(driver driver.Bucket) *Bucket {
	return &Bucket{b: driver}
}
