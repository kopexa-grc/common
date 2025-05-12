// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package driver

import (
	"context"
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
