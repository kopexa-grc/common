// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package azurestore

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/kopexa-grc/common/blob/driver"
)

type AzureStore struct {
	Service   AzService
	Container string
}

func New(service AzService) *AzureStore {
	return &AzureStore{
		Service: service,
	}
}

func (store *AzureStore) SignedURL(ctx context.Context, key string, opts *driver.SignedURLOptions) (string, error) {
	blob, err := store.Service.NewBlob(ctx, key)
	if err != nil {
		return "", err
	}

	return blob.SignedURL(ctx, opts)
}

func (store *AzureStore) GetSignedUploadURL(ctx context.Context, key string, expires time.Duration, _ int64, contentType string) (string, error) {
	blob, err := store.Service.NewBlob(ctx, key)
	if err != nil {
		return "", err
	}

	return blob.SignedURL(ctx, &driver.SignedURLOptions{
		Expiry:      expires,
		Method:      http.MethodPut,
		ContentType: contentType,
	})
}

func (store *AzureStore) NewRangeReader(ctx context.Context, key string, offset, length int64, opts *driver.ReaderOptions) (driver.Reader, error) {
	blob, err := store.Service.NewBlob(ctx, key)
	if err != nil {
		return nil, err
	}

	return blob.NewRangeReader(ctx, key, offset, length, opts)
}

func (store *AzureStore) NewTypedWriter(ctx context.Context, key, contentType string, opts *driver.WriterOptions) (driver.Writer, error) {
	blob, err := store.Service.NewBlob(ctx, key)
	if err != nil {
		return nil, err
	}

	return blob.NewTypedWriter(ctx, key, contentType, opts)
}

func (store *AzureStore) GetSignedDownloadURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	blob, err := store.Service.NewBlob(ctx, key)
	if err != nil {
		return "", err
	}

	return blob.SignedURL(ctx, &driver.SignedURLOptions{
		Expiry: expires,
		Method: http.MethodGet,
	})
}

// DeleteObject is a wrapper around the Delete method for
// compatibility with the StorageProvider interface.
func (store *AzureStore) DeleteObject(ctx context.Context, key string) error {
	return store.Delete(ctx, key)
}

func (store *AzureStore) Delete(ctx context.Context, key string) error {
	blob, err := store.Service.NewBlob(ctx, key)
	if err != nil {
		return err
	}

	return blob.Delete(ctx)
}

func (store *AzureStore) Copy(ctx context.Context, dstKey, srcKey string, opts *driver.CopyOptions) error {
	dstBlobClient, err := store.Service.NewBlob(ctx, dstKey)
	if err != nil {
		return err
	}

	srcBlobClient, err := store.Service.NewBlob(ctx, srcKey)
	if err != nil {
		return err
	}

	resp, err := dstBlobClient.StartCopyFromURL(ctx, srcBlobClient.URL(), opts)
	if err != nil {
		return err
	}

	nErrors := 0

	copyStatus := *resp.CopyStatus
	for copyStatus == blob.CopyStatusTypePending {
		time.Sleep(defaultCopyPollMs * time.Millisecond)

		propertiesResp, err := dstBlobClient.GetProperties(ctx, nil)
		if err != nil {
			nErrors++
			if ctx.Err() != nil || nErrors == 3 {
				return err
			}

			continue
		}

		copyStatus = *propertiesResp.CopyStatus
	}

	if copyStatus != blob.CopyStatusTypeSuccess {
		return fmt.Errorf("%w: %s", driver.ErrCopyFailed, copyStatus)
	}

	return nil
}

func (store *AzureStore) TestConnection() error {
	return nil
}
