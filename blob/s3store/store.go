// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package s3store

import (
	"context"
	"io"
	"time"

	"github.com/kopexa-grc/common/blob/driver"
)

type Store struct {
	Service S3Service
}

func New(service *S3Service) *Store {
	if service == nil {
		return nil
	}

	return &Store{
		Service: *service,
	}
}

func (s *Store) Delete(ctx context.Context, key string) error {
	return s.Service.Delete(ctx, key)
}

func (s *Store) SignedURL(ctx context.Context, key string, opts *driver.SignedURLOptions) (string, error) {
	exp := 15 * time.Minute
	if opts != nil && opts.Expiry > 0 {
		exp = opts.Expiry
	}

	method := "GET"
	if opts != nil && opts.Method != "" {
		method = opts.Method
	}

	return s.Service.GetSignedURL(ctx, key, exp, method)
}

func (s *Store) Copy(ctx context.Context, srcKey, dstKey string, opts *driver.CopyOptions) error {
	return s.Service.Copy(ctx, CopyParams{
		SourceKey: srcKey,
		DestKey:   dstKey,
	})
}

func (s *Store) NewTypedWriter(ctx context.Context, key, contentType string, opts *driver.WriterOptions) (driver.Writer, error) {
	pr, pw := io.Pipe()
	errCh := make(chan error, 1)

	go func() {
		err := s.Service.UploadTyped(ctx, key, pr, contentType, opts)
		errCh <- err
		close(errCh)
	}()

	return &writer{
		pw:    pw,
		errCh: errCh,
	}, nil
}

func (s *Store) NewRangeReader(ctx context.Context, key string, offset, length int64, opts *driver.ReaderOptions) (driver.Reader, error) {
	r, err := s.Service.RangeDownload(ctx, key, offset, length)
	if err != nil {
		return nil, err
	}

	attrs := &driver.ReaderAttributes{
		ContentType: r.ContentType(),
		Size:        r.Size(),
		ModTime:     r.ModTime(),
	}

	return &reader{
		ReadCloser: r,
		attrs:      attrs,
	}, nil
}

// -- helper types --

type reader struct {
	io.ReadCloser
	body  io.ReadCloser
	attrs *driver.ReaderAttributes
}

func (r *reader) Attributes() *driver.ReaderAttributes { return r.attrs }

func (r *reader) As(i any) bool {
	if v, ok := i.(*io.ReadCloser); ok {
		*v = r.ReadCloser
		return true
	}
	return false
}

type writer struct {
	pw    *io.PipeWriter
	errCh chan error
}

func (w *writer) Write(p []byte) (int, error) { return w.pw.Write(p) }

func (w *writer) Close() error {
	if err := w.pw.Close(); err != nil {
		return err
	}
	return <-w.errCh
}
