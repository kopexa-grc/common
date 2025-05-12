// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package azurestore

import (
	"context"
	"fmt"
	"net/http"
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
	"github.com/rs/zerolog/log"
)

type AzBlob interface {
	SignedURL(ctx context.Context, opts *driver.SignedURLOptions) (string, error)
	StartCopyFromURL(ctx context.Context, url string, opts *driver.CopyOptions) (blob.StartCopyFromURLResponse, error)
	GetProperties(ctx context.Context, o *blob.GetPropertiesOptions) (blob.GetPropertiesResponse, error)
	Delete(ctx context.Context) error
	URL() string
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
