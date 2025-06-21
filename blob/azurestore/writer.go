// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package azurestore

import (
	"context"
	"io"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
)

type writer struct {
	ctx        context.Context
	client     *blockblob.Client
	uploadOpts *azblob.UploadStreamOptions

	// Ends of an io.Pipe, created when the first byte is written.
	pw *io.PipeWriter
	pr *io.PipeReader

	// Alternatively, upload is set to true when Upload was
	// used to upload data.
	upload bool

	donec chan struct{} // closed when done writing
	// The following fields will be written before donec closes:
	err error
}

// Write appends p to w.pw. User must call Close to close the w after done writing.
func (w *writer) Write(p []byte) (int, error) {
	// Avoid opening the pipe for a zero-length write;
	// the concrete can do these for empty blobs.
	if len(p) == 0 {
		return 0, nil
	}
	if w.pw == nil {
		// We'll write into pw and use pr as an io.Reader for the
		// Upload call to Azure.
		w.pr, w.pw = io.Pipe()
		w.open(w.pr, true)
	}
	return w.pw.Write(p)
}

// Upload reads from r. Per the driver, it is guaranteed to be the only
// write call for this writer.
func (w *writer) Upload(r io.Reader) error {
	w.upload = true
	w.open(r, false)
	return nil
}

// r may be nil if we're Closing and no data was written.
// If closePipeOnError is true, w.pr will be closed if there's an
// error uploading to Azure.
func (w *writer) open(r io.Reader, closePipeOnError bool) {
	go func() {
		defer close(w.donec)

		if r == nil {
			r = http.NoBody
		}
		_, w.err = w.client.UploadStream(w.ctx, r, w.uploadOpts)
		if w.err != nil {
			if closePipeOnError {
				w.pr.CloseWithError(w.err)
				w.pr = nil
			}
		}
	}()
}

// Close completes the writer and closes it. Any error occurring during write will
// be returned. If a writer is closed before any Write is called, Close will
// create an empty file at the given key.
func (w *writer) Close() error {
	if !w.upload {
		if w.pr != nil {
			defer w.pr.Close()
		}
		if w.pw == nil {
			// We never got any bytes written. We'll write an http.NoBody.
			w.open(nil, false)
		} else if err := w.pw.Close(); err != nil {
			return err
		}
	}
	<-w.donec
	return w.err
}
