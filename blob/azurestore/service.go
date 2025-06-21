// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package azurestore

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/kopexa-grc/common/blob/driver"
	"github.com/kopexa-grc/common/blob/internal/escape"
	kerr "github.com/kopexa-grc/common/errors"
	"github.com/rs/zerolog/log"
)

type AzBlob interface {
	SignedURL(ctx context.Context, opts *driver.SignedURLOptions) (string, error)
	StartCopyFromURL(ctx context.Context, url string, opts *driver.CopyOptions) (blob.StartCopyFromURLResponse, error)
	GetProperties(ctx context.Context, o *blob.GetPropertiesOptions) (blob.GetPropertiesResponse, error)
	Delete(ctx context.Context) error
	URL() string
	NewRangeReader(ctx context.Context, offset, length int64, opts *driver.ReaderOptions) (driver.Reader, error)
	NewTypedWriter(ctx context.Context, contentType string, opts *driver.WriterOptions) (driver.Writer, error)
}

type BlockBlob struct {
	BlobClient     *blockblob.Client
	Indexes        []int
	BlobAccessTier *blob.AccessTier
}

type AzService interface {
	NewBlob(ctx context.Context, name string) (AzBlob, error)
}

type azService struct {
	ContainerClient *container.Client
	ContainerName   string
	BlobAccessTier  *blob.AccessTier
}

type AzConfig struct {
	AccountName         string
	AccountKey          string
	BlobAccessTier      string
	ContainerName       string
	ContainerAccessType string
	Endpoint            string
}

const (
	defaultMaxRetries    = 5
	defaultRetryDelay    = 100  // ms
	defaultMaxRetryDelay = 5000 // ms
	defaultCopyPollMs    = 500  // ms
)

func NewAzureService(config *AzConfig) (AzService, error) {
	cred, err := azblob.NewSharedKeyCredential(config.AccountName, config.AccountKey)
	if err != nil {
		return nil, err
	}

	serviceURL := fmt.Sprintf("%s/%s", config.Endpoint, config.ContainerName)
	retryOpts := policy.RetryOptions{
		MaxRetries:    maxRetries,
		RetryDelay:    retryDelay,    // Retry after 100ms initially
		MaxRetryDelay: maxRetryDelay, // Max retry delay 5 seconds
	}

	containerClient, err := container.NewClientWithSharedKeyCredential(serviceURL, cred, &container.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Retry: retryOpts,
		},
	})
	if err != nil {
		return nil, err
	}

	containerCreateOptions := &container.CreateOptions{}

	switch config.ContainerAccessType {
	case "container":
		containerCreateOptions.Access = to.Ptr(container.PublicAccessTypeContainer)
	case "blob":
		containerCreateOptions.Access = to.Ptr(container.PublicAccessTypeBlob)
	default:
		// Leaving Access nil will default to private access
	}

	_, err = containerClient.Create(context.Background(), containerCreateOptions)
	//nolint:gocritic
	if err != nil && !strings.Contains(err.Error(), "ContainerAlreadyExists") {
		return nil, err
	} else if err == nil {
		log.Info().Str("container", config.ContainerName).Msg("Azure Blob container created")
	} else {
		log.Debug().Str("container", config.ContainerName).Msg("Azure Blob container already exists")
	}

	var blobAccessTier *blob.AccessTier

	switch config.BlobAccessTier {
	case "archive":
		blobAccessTier = to.Ptr(blob.AccessTierArchive)
	case "cool":
		blobAccessTier = to.Ptr(blob.AccessTierCool)
	case "hot":
		blobAccessTier = to.Ptr(blob.AccessTierHot)
	}

	return &azService{
		ContainerClient: containerClient,
		ContainerName:   config.ContainerName,
		BlobAccessTier:  blobAccessTier,
	}, nil
}

// Determine if we return a InfoBlob or BlockBlob, based on the name
func (service *azService) NewBlob(_ context.Context, name string) (AzBlob, error) {
	blobClient := service.ContainerClient.NewBlockBlobClient(escapeKey(name, false))

	return &BlockBlob{
		BlobClient:     blobClient,
		Indexes:        []int{},
		BlobAccessTier: service.BlobAccessTier,
	}, nil
}

func (blockBlob *BlockBlob) SignedURL(_ context.Context, opts *driver.SignedURLOptions) (string, error) {
	perms := sas.BlobPermissions{}

	switch opts.Method {
	case http.MethodGet:
		perms.Read = true
	case http.MethodPut:
		perms.Create = true
		perms.Write = true
	case http.MethodDelete:
		perms.Delete = true
	default:
		return "", driver.ErrUnsupportedMethod
	}

	if opts.BeforeSign != nil {
		asFunc := func(i any) bool {
			v, ok := i.(**sas.BlobPermissions)
			if ok {
				*v = &perms
			}

			return ok
		}
		if err := opts.BeforeSign(asFunc); err != nil {
			return "", err
		}
	}

	start := time.Now().UTC()
	expiry := start.Add(opts.Expiry)

	return blockBlob.BlobClient.GetSASURL(perms, expiry, &blob.GetSASURLOptions{StartTime: &start})
}

// Delete the blockBlob from Azure Blob Storage
func (blockBlob *BlockBlob) Delete(ctx context.Context) error {
	deleteOptions := &azblob.DeleteBlobOptions{
		DeleteSnapshots: to.Ptr(azblob.DeleteSnapshotsOptionTypeInclude),
	}
	_, err := blockBlob.BlobClient.Delete(ctx, deleteOptions)

	return err
}

// StartCopyFromURL starts a copy operation from a URL to the blockBlob
func (blockBlob *BlockBlob) StartCopyFromURL(ctx context.Context, url string, opts *driver.CopyOptions) (blob.StartCopyFromURLResponse, error) {
	copyOptions := &blob.StartCopyFromURLOptions{}

	if opts.BeforeCopy != nil {
		asFunc := func(i any) bool {
			//nolint:gocritic
			switch v := i.(type) {
			case **blob.StartCopyFromURLOptions:
				*v = copyOptions
				return true
			}

			return false
		}
		if err := opts.BeforeCopy(asFunc); err != nil {
			return blob.StartCopyFromURLResponse{}, err
		}
	}

	return blockBlob.BlobClient.StartCopyFromURL(ctx, url, copyOptions)
}

// URL returns the URL of the blockBlob
func (blockBlob *BlockBlob) URL() string {
	return blockBlob.BlobClient.URL()
}

// GetProperties gets the properties of the blockBlob
func (blockBlob *BlockBlob) GetProperties(ctx context.Context, o *blob.GetPropertiesOptions) (blob.GetPropertiesResponse, error) {
	return blockBlob.BlobClient.GetProperties(ctx, o)
}

// reader reads an azblob. It implements io.ReadCloser.
type reader struct {
	body  io.ReadCloser
	attrs driver.ReaderAttributes
	raw   *azblob.DownloadStreamResponse
}

func (r *reader) Read(p []byte) (int, error) {
	return r.body.Read(p)
}

func (r *reader) Close() error {
	return r.body.Close()
}

func (r *reader) Attributes() *driver.ReaderAttributes {
	return &r.attrs
}

func (r *reader) As(i any) bool {
	p, ok := i.(*azblob.DownloadStreamResponse)
	if !ok {
		return false
	}

	*p = *r.raw

	return true
}

// NewRangeReader implements driver.NewRangeReader.
func (blockBlob *BlockBlob) NewRangeReader(ctx context.Context, offset, length int64, opts *driver.ReaderOptions) (driver.Reader, error) {
	blobClient := blockBlob.BlobClient

	downloadOpts := azblob.DownloadStreamOptions{}
	if offset != 0 {
		downloadOpts.Range.Offset = offset
	}

	if length >= 0 {
		downloadOpts.Range.Count = length
	}

	if opts.BeforeRead != nil {
		asFunc := func(i any) bool {
			if p, ok := i.(**azblob.DownloadStreamOptions); ok {
				*p = &downloadOpts
				return true
			}

			return false
		}
		if err := opts.BeforeRead(asFunc); err != nil {
			return nil, err
		}
	}

	blobDownloadResponse, err := blobClient.DownloadStream(ctx, &downloadOpts)
	if err != nil {
		return nil, err
	}

	attrs := driver.ReaderAttributes{
		ContentType: *blobDownloadResponse.ContentType,
		Size:        getSize(blobDownloadResponse.ContentLength, *blobDownloadResponse.ContentRange),
		ModTime:     *blobDownloadResponse.LastModified,
	}

	var body io.ReadCloser
	if length == 0 {
		body = http.NoBody
	} else {
		body = blobDownloadResponse.Body
	}

	return &reader{
		body:  body,
		attrs: attrs,
		raw:   &blobDownloadResponse,
	}, nil
}

func (blockBlob *BlockBlob) NewTypedWriter(ctx context.Context, contentType string, opts *driver.WriterOptions) (driver.Writer, error) {
	blobClient := blockBlob.BlobClient

	if opts.BufferSize == 0 {
		opts.BufferSize = defaultUploadBlockSize
	}

	if opts.MaxConcurrency == 0 {
		opts.MaxConcurrency = defaultUploadBuffers
	}

	md := make(map[string]*string, len(opts.Metadata))

	for k, v := range opts.Metadata {
		// See the package comments for more details on escaping of metadata
		// keys & values.
		e := escape.HexEscape(k, func(runes []rune, i int) bool {
			c := runes[i]

			switch {
			case i == 0 && c >= '0' && c <= '9':
				return true
			case escape.IsASCIIAlphanumeric(c):
				return false
			case c == '_':
				return false
			}

			return true
		})
		if _, ok := md[e]; ok {
			return nil, kerr.Newf(kerr.InvalidArgument, nil, "duplicate keys after escaping: %q => %q", k, e)
		}

		escaped := escape.URLEscape(v)
		md[e] = &escaped
	}

	uploadOpts := &azblob.UploadStreamOptions{
		BlockSize:   int64(opts.BufferSize),
		Concurrency: opts.MaxConcurrency,
		Metadata:    md,
		HTTPHeaders: &blob.HTTPHeaders{
			BlobCacheControl:       &opts.CacheControl,
			BlobContentDisposition: &opts.ContentDisposition,
			BlobContentEncoding:    &opts.ContentEncoding,
			BlobContentLanguage:    &opts.ContentLanguage,
			BlobContentMD5:         opts.ContentMD5,
			BlobContentType:        &contentType,
		},
	}

	if opts.IfNotExist {
		etagAny := azcore.ETagAny
		uploadOpts.AccessConditions = &azblob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &etagAny,
			},
		}
	}

	if opts.BeforeWrite != nil {
		asFunc := func(i any) bool {
			p, ok := i.(**azblob.UploadStreamOptions)
			if !ok {
				return false
			}

			*p = uploadOpts

			return true
		}
		if err := opts.BeforeWrite(asFunc); err != nil {
			return nil, err
		}
	}

	return &writer{
		ctx:        ctx,
		client:     blobClient,
		uploadOpts: uploadOpts,
		donec:      make(chan struct{}),
	}, nil
}

func getSize(contentLength *int64, contentRange string) int64 {
	var size int64
	// Default size to ContentLength, but that's incorrect for partial-length reads,
	// where ContentLength refers to the size of the returned Body, not the entire
	// size of the blob. ContentRange has the full size.
	if contentLength != nil {
		size = *contentLength
	}

	if contentRange != "" {
		// Sample: bytes 10-14/27 (where 27 is the full size).
		parts := strings.Split(contentRange, "/")

		const expectedParts = 2

		if len(parts) == expectedParts {
			if i, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				size = i
			}
		}
	}

	return size
}

// escapeKey does all required escaping for UTF-8 strings to work
// with Azure. isPrefix indicated whether this is a prefix/delimeter or the full key.
func escapeKey(key string, isPrefix bool) string {
	return escape.HexEscape(key, func(r []rune, i int) bool {
		c := r[i]

		switch {
		// Azure does not work well with backslashes in blob names.
		case c == '\\':
			return true
		// Azure doesn't handle these characters (determined via experimentation).
		case c < 32 || c == 34 || c == 35 || c == 37 || c == 63 || c == 127:
			return true
		// Escape trailing "/" for full keys, otherwise Azure can't address them
		// consistently.
		case !isPrefix && i == len(key)-1 && c == '/':
			return true
		// For "../", escape the trailing slash.
		case i > 1 && r[i] == '/' && r[i-1] == '.' && r[i-2] == '.':
			return true
		}

		return false
	})
}
