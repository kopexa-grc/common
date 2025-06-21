// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob

import (
	"bytes"
	"context"
	"hash"
	"io"
	"net/http"

	"github.com/kopexa-grc/common/blob/driver"
	kerr "github.com/kopexa-grc/common/errors"
	"go.opentelemetry.io/otel/metric"
)

// Writer writes bytes to a blob.
//
// It implements io.WriteCloser (https://golang.org/pkg/io/#Closer), and must be
// closed after all writes are done.
type Writer struct {
	b          driver.Bucket
	w          driver.Writer
	key        string
	end        func(err error) // called at Close to finish trace and metric collection
	cancel     func()          // cancels the ctx provided to NewTypedWriter if contentMD5 verification fails
	contentMD5 []byte
	md5hash    hash.Hash

	// Metric collection fields
	bytesWrittenCounter metric.Int64Counter
	bytesWritten        int
	closed              bool

	// These fields are non-zero values only when w is nil (not yet created).
	//
	// A ctx is stored in the Writer since we need to pass it into NewTypedWriter
	// when we finish detecting the content type of the blob and create the
	// underlying driver.Writer. This step happens inside Write or Close and
	// neither of them take a context.Context as an argument.
	//
	// All 3 fields are only initialized when we create the Writer without
	// setting the w field, and are reset to zero values after w is created.
	ctx  context.Context
	opts *driver.WriterOptions
	buf  *bytes.Buffer
}

// sniffLen is the byte size of Writer.buf used to detect content-type.
const sniffLen = 512

// Write implements the io.Writer interface (https://golang.org/pkg/io/#Writer).
//
// Writes may happen asynchronously, so the returned error can be nil
// even if the actual write eventually fails. The write is only guaranteed to
// have succeeded if Close returns no error.
func (w *Writer) Write(p []byte) (int, error) {
	if len(w.contentMD5) > 0 {
		if _, err := w.md5hash.Write(p); err != nil {
			return 0, err
		}
	}

	if w.w != nil {
		return w.write(p)
	}

	// If w is not yet created due to no content-type being passed in, try to sniff
	// the MIME type based on at most 512 bytes of the blob content of p.

	// Detect the content-type directly if the first chunk is at least 512 bytes.
	if w.buf.Len() == 0 && len(p) >= sniffLen {
		return w.open(p)
	}

	// Store p in w.buf and detect the content-type when the size of content in
	// w.buf is at least 512 bytes.
	n, err := w.buf.Write(p)
	if err != nil {
		return 0, err
	}

	if w.buf.Len() >= sniffLen {
		// Note that w.open will return the full length of the buffer; we don't want
		// to return that as the length of this write since some of them were written in
		// previous writes. Instead, we return the n from this write, above.
		_, err := w.open(w.buf.Bytes())
		return n, err
	}

	return n, nil
}

// Close closes the blob writer. The write operation is not guaranteed to have succeeded until
// Close returns with no error.
// Close may return an error if the context provided to create the Writer is
// canceled or reaches its deadline.
func (w *Writer) Close() (err error) {
	w.closed = true

	// Store context before it might be set to nil in open()
	ctx := w.ctx

	defer func() {
		if w.end != nil {
			w.end(err)
		}
		// Emit only on close to avoid an allocation on each call to Write().
		// Record bytes written metric with OpenTelemetry
		if w.bytesWrittenCounter != nil && w.bytesWritten > 0 && ctx != nil {
			w.bytesWrittenCounter.Add(
				ctx,
				int64(w.bytesWritten))
		}
	}()

	if len(w.contentMD5) > 0 {
		// Verify the MD5 hash of what was written matches the ContentMD5 provided
		// by the user.
		md5sum := w.md5hash.Sum(nil)
		if !bytes.Equal(md5sum, w.contentMD5) {
			// No match! Return an error, but first cancel the context and call the
			// driver's Close function to ensure the write is aborted.
			w.cancel()

			if w.w != nil {
				_ = w.w.Close()
			}

			return kerr.Newf(kerr.FailedPrecondition, nil, "blob: the WriterOptions.ContentMD5 you specified (%X) did not match what was written (%X)", w.contentMD5, md5sum)
		}
	}

	defer w.cancel()

	if w.w != nil {
		return wrapError(w.b, w.w.Close(), w.key)
	}

	if _, err := w.open(w.buf.Bytes()); err != nil {
		return err
	}

	return wrapError(w.b, w.w.Close(), w.key)
}

// open tries to detect the MIME type of p and write it to the blob.
// The error it returns is wrapped.
func (w *Writer) open(p []byte) (int, error) {
	ct := http.DetectContentType(p)

	var err error

	if w.w, err = w.b.NewTypedWriter(w.ctx, w.key, ct, w.opts); err != nil {
		return 0, wrapError(w.b, err, w.key)
	}
	// Set the 3 fields needed for lazy NewTypedWriter back to zero values
	// (see the comment on Writer).
	w.buf = nil
	w.ctx = nil
	w.opts = nil

	return w.write(p)
}

func (w *Writer) write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.bytesWritten += n

	return n, wrapError(w.b, err, w.key)
}

// ReadFrom reads from r and writes to w until EOF or error.
// The return value is the number of bytes read from r.
//
// It implements the io.ReaderFrom interface.
func (w *Writer) ReadFrom(r io.Reader) (int64, error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Don't do this for our own *Reader to avoid infinite recursion.
	// Avoids an allocation and a copy.
	switch r.(type) {
	case *Reader:
	default:
		if wt, ok := r.(io.WriterTo); ok {
			n, err := wt.WriteTo(w)
			return n, err
		}
	}

	nr, _, err := readFromWriteTo(r, w)

	return nr, err
}

// uploadAndClose is similar to ReadFrom, but ensures it's the only write.
// This pattern is more optimal for some drivers.
func (w *Writer) uploadAndClose(r io.Reader) (err error) {
	if w.bytesWritten != 0 {
		// Shouldn't happen.
		return kerr.Newf(kerr.UnexpectedFailure, nil, "blob: uploadAndClose must be the first write")
	}
	// When ContentMD5 is being checked, we can't use Upload.
	if len(w.contentMD5) > 0 {
		_, err = w.ReadFrom(r)
	} else {
		driverUploader, ok := w.w.(driver.Uploader)
		if ok {
			err = driverUploader.Upload(r)
		} else {
			_, err = w.ReadFrom(r)
		}
	}

	cerr := w.Close()
	if err == nil && cerr != nil {
		err = cerr
	}

	return err
}
