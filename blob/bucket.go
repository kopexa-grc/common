// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"runtime"
	"strings"
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

// ReaderOptions sets options for NewReader and NewRangeReader.
type ReaderOptions struct {
	// BeforeRead is a callback that will be called before
	// any data is read (unless NewReader returns an error before then, in which
	// case it may not be called at all).
	//
	// Calling Seek may reset the underlying reader, and result in BeforeRead
	// getting called again with a different underlying provider-specific reader..
	//
	// asFunc converts its argument to driver-specific types.
	// See https://gocloud.dev/concepts/as/ for background information.
	BeforeRead func(asFunc func(any) bool) error
}

// WriterOptions sets options for NewWriter.
type WriterOptions struct {
	// BufferSize changes the default size in bytes of the chunks that
	// Writer will upload in a single request; larger blobs will be split into
	// multiple requests.
	//
	// This option may be ignored by some drivers.
	//
	// If 0, the driver will choose a reasonable default.
	//
	// If the Writer is used to do many small writes concurrently, using a
	// smaller BufferSize may reduce memory usage.
	BufferSize int

	// MaxConcurrency changes the default concurrency for parts of an upload.
	//
	// This option may be ignored by some drivers.
	//
	// If 0, the driver will choose a reasonable default.
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

	// ContentType specifies the MIME type of the blob being written. If not set,
	// it will be inferred from the content using the algorithm described at
	// http://mimesniff.spec.whatwg.org/.
	// Set DisableContentTypeDetection to true to disable the above and force
	// the ContentType to stay empty.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
	ContentType string

	// When true, if ContentType is the empty string, it will stay the empty
	// string rather than being inferred from the content.
	// Note that while the blob will be written with an empty string ContentType,
	// most providers will fill one in during reads, so don't expect an empty
	// ContentType if you read the blob back.
	DisableContentTypeDetection bool

	// ContentMD5 is used as a message integrity check.
	// If len(ContentMD5) > 0, the MD5 hash of the bytes written must match
	// ContentMD5, or Close will return an error without completing the write.
	// https://tools.ietf.org/html/rfc1864
	ContentMD5 []byte

	// Metadata holds key/value strings to be associated with the blob, or nil.
	// Keys may not be empty, and are lowercased before being written.
	// Duplicate case-insensitive keys (e.g., "foo" and "FOO") will result in
	// an error.
	Metadata map[string]string

	// BeforeWrite is a callback that will be called exactly once, before
	// any data is written (unless NewWriter returns an error, in which case
	// it will not be called at all). Note that this is not necessarily during
	// or after the first Write call, as drivers may buffer bytes before
	// sending an upload request.
	//
	// asFunc converts its argument to driver-specific types.
	// See https://gocloud.dev/concepts/as/ for background information.
	BeforeWrite func(asFunc func(any) bool) error

	// IfNotExist is used for conditional writes. When set to 'true',
	// if a blob exists for the same key in the bucket, the write
	// operation won't succeed and the current blob for the key will
	// be left untouched. An error for which gcerrors.Code will return
	// gcerrors.PreconditionFailed will be returned by Write or Close.
	IfNotExist bool
}

// Uploads reads from a io.Reader and writes into a blob
//
// opts.ContentType is required.
func (b *Bucket) Upload(ctx context.Context, key string, r io.Reader, opts *WriterOptions) (err error) {
	if opts == nil || opts.ContentType == "" {
		return kerr.Newf(kerr.InvalidArgument, nil, "blob: Upload requires WriterOptions.ContentType")
	}
	w, err := b.NewWriter(ctx, key, opts)
	if err != nil {
		return err
	}
	return w.uploadAndClose(r)
}

// NewWriter returns a Writer that writes to the blob stored at key.
// A nil WriterOptions is treated the same as the zero value.
//
// If a blob with this key already exists, it will be replaced.
// The blob being written is not guaranteed to be readable until Close
// has been called; until then, any previous blob will still be readable.
// Even after Close is called, newly written blobs are not guaranteed to be
// returned from List; some services are only eventually consistent.
//
// The returned Writer will store ctx for later use in Write and/or Close.
// To abort a write, cancel ctx; otherwise, it must remain open until
// Close is called.
//
// The caller must call Close on the returned Writer, even if the write is
// aborted.
func (b *Bucket) NewWriter(ctx context.Context, key string, opts *WriterOptions) (_ *Writer, err error) {
	if !utf8.ValidString(key) {
		return nil, kerr.Newf(kerr.InvalidArgument, nil, "blob: NewWriter key must be a valid UTF-8 string: %q", key)
	}
	if opts == nil {
		opts = &WriterOptions{}
	}
	dopts := &driver.WriterOptions{
		CacheControl:                opts.CacheControl,
		ContentDisposition:          opts.ContentDisposition,
		ContentEncoding:             opts.ContentEncoding,
		ContentLanguage:             opts.ContentLanguage,
		ContentMD5:                  opts.ContentMD5,
		BufferSize:                  opts.BufferSize,
		MaxConcurrency:              opts.MaxConcurrency,
		BeforeWrite:                 opts.BeforeWrite,
		DisableContentTypeDetection: opts.DisableContentTypeDetection,
		IfNotExist:                  opts.IfNotExist,
	}
	if len(opts.Metadata) > 0 {
		// Services are inconsistent, but at least some treat keys
		// as case-insensitive. To make the behavior consistent, we
		// force-lowercase them when writing and reading.
		md := make(map[string]string, len(opts.Metadata))
		for k, v := range opts.Metadata {
			if k == "" {
				return nil, kerr.Newf(kerr.InvalidArgument, nil, "blob: WriterOptions.Metadata keys may not be empty strings")
			}
			if !utf8.ValidString(k) {
				return nil, kerr.Newf(kerr.InvalidArgument, nil, "blob: WriterOptions.Metadata keys must be valid UTF-8 strings: %q", k)
			}
			if !utf8.ValidString(v) {
				return nil, kerr.Newf(kerr.InvalidArgument, nil, "blob: WriterOptions.Metadata values must be valid UTF-8 strings: %q", v)
			}
			lowerK := strings.ToLower(k)
			if _, found := md[lowerK]; found {
				return nil, kerr.Newf(kerr.InvalidArgument, nil, "blob: WriterOptions.Metadata has a duplicate case-insensitive metadata key: %q", lowerK)
			}
			md[lowerK] = v
		}
		dopts.Metadata = md
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, errClosed
	}
	ctx, cancel := context.WithCancel(ctx)

	w := &Writer{
		b:          b.b,
		cancel:     cancel,
		key:        key,
		contentMD5: opts.ContentMD5,
		md5hash:    md5.New(),
		ctx:        ctx,
	}
	if opts.ContentType != "" || opts.DisableContentTypeDetection {
		var ct string
		if opts.ContentType != "" {
			t, p, err := mime.ParseMediaType(opts.ContentType)
			if err != nil {
				cancel()
				return nil, err
			}
			ct = mime.FormatMediaType(t, p)
		}
		dw, err := b.b.NewTypedWriter(ctx, key, ct, dopts)
		if err != nil {
			cancel()
			return nil, wrapError(b.b, err, key)
		}
		w.w = dw
	} else {
		// Save the fields needed to called NewTypedWriter later, once we've gotten
		// sniffLen bytes; see the comment on Writer.
		w.opts = dopts
		w.buf = bytes.NewBuffer([]byte{})
	}
	_, file, lineno, ok := runtime.Caller(1)
	runtime.SetFinalizer(w, func(w *Writer) {
		if !w.closed {
			var caller string
			if ok {
				caller = fmt.Sprintf(" (%s:%d)", file, lineno)
			}
			log.Printf("A blob.Writer writing to %q was never closed%s", key, caller)
		}
	})
	return w, nil
}

// NewRangeReader returns a Reader to read content from the blob stored at key.
// It reads at most length bytes starting at offset (>= 0).
// If length is negative, it will read till the end of the blob.
//
// For the purposes of Seek, the returned Reader will start at offset and
// end at the minimum of the actual end of the blob or (if length > 0) offset + length.
//
// Note that ctx is used for all reads performed during the lifetime of the reader.
//
// If the blob does not exist, NewRangeReader returns an error for which
// gcerrors.Code will return gcerrors.NotFound. Exists is a lighter-weight way
// to check for existence.
//
// A nil ReaderOptions is treated the same as the zero value.
//
// The caller must call Close on the returned Reader when done reading.
func (b *Bucket) NewRangeReader(ctx context.Context, key string, offset, length int64, opts *ReaderOptions) (_ *Reader, err error) {
	return b.newRangeReader(ctx, key, offset, length, opts)
}

func (b *Bucket) newRangeReader(ctx context.Context, key string, offset, length int64, opts *ReaderOptions) (_ *Reader, err error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, errClosed
	}
	if offset < 0 {
		return nil, kerr.Newf(kerr.InvalidArgument, nil, "blob: NewRangeReader offset must be non-negative (%d)", offset)
	}
	if !utf8.ValidString(key) {
		return nil, kerr.Newf(kerr.InvalidArgument, nil, "blob: NewRangeReader key must be a valid UTF-8 string: %q", key)
	}
	if opts == nil {
		opts = &ReaderOptions{}
	}
	dopts := &driver.ReaderOptions{
		BeforeRead: opts.BeforeRead,
	}

	var dr driver.Reader
	dr, err = b.b.NewRangeReader(ctx, key, offset, length, dopts)
	if err != nil {
		return nil, wrapError(b.b, err, key)
	}

	r := &Reader{
		b:           b.b,
		r:           dr,
		key:         key,
		ctx:         ctx,
		dopts:       dopts,
		baseOffset:  offset,
		baseLength:  length,
		savedOffset: -1,
	}
	_, file, lineno, ok := runtime.Caller(2)
	runtime.SetFinalizer(r, func(r *Reader) {
		if !r.closed {
			var caller string
			if ok {
				caller = fmt.Sprintf(" (%s:%d)", file, lineno)
			}
			log.Printf("A blob.Reader reading from %q was never closed%s", key, caller)
		}
	})
	return r, nil
}

// NewBucketForTest creates a Bucket with a mock driver for testing purposes.
func NewBucketForTest(driver driver.Bucket) *Bucket {
	return &Bucket{b: driver}
}
