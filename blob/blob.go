// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package blob provides a portable interface for blob storage operations.
//
// This package implements the Google Cloud Storage Go client library pattern
// for blob storage, providing a unified interface across different storage
// providers. The package supports Azure Blob Storage as the primary backend
// with plans for additional provider support.
//
// The blob package follows the Google API Design Guide principles:
//   - Resource-oriented design with clear resource naming
//   - Standard HTTP methods for operations (GET, PUT, DELETE)
//   - Consistent error handling with structured error responses
//   - Support for conditional operations and optimistic concurrency
//   - Comprehensive metadata and content type handling
//
// Example usage:
//
//	config := &blob.Config{
//		Azure: blob.AzureConfig{
//			AccountName: "myaccount",
//			AccountKey:  "mykey",
//			Endpoint:    "https://myaccount.blob.core.windows.net",
//		},
//	}
//
//	provider, err := blob.New(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Access public bucket
//	publicBucket, err := provider.Public()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Upload a file
//	file, _ := os.Open("example.txt")
//	defer file.Close()
//	err = publicBucket.Upload(ctx, "example.txt", file, nil)
//
//	// Access space-specific bucket
//	spaceBucket, err := provider.Space("workspace-123")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The package provides two types of buckets:
//   - Public buckets: Shared storage accessible to all users
//   - Space buckets: Isolated storage for specific workspaces or projects
//
// All operations are context-aware and support cancellation, timeouts,
// and request tracing. The package also provides support for signed URLs,
// range reads, and conditional operations.
package blob

import (
	"fmt"

	"github.com/kopexa-grc/common/blob/azurestore"
)

// Config represents the configuration for blob storage operations.
//
// The configuration supports multiple storage providers, with Azure Blob Storage
// being the primary supported backend. Additional providers can be added by
// extending this configuration structure.
//
// The configuration follows the Google API Design Guide principle of using
// structured configuration objects rather than individual parameters.
type Config struct {
	// Azure contains the configuration for Azure Blob Storage.
	// This is the primary supported storage backend.
	Azure AzureConfig
}

// AzureConfig contains the configuration parameters for Azure Blob Storage.
//
// Azure Blob Storage requires an account name, authentication key, and endpoint
// URL. The endpoint typically follows the pattern:
// https://{account-name}.blob.core.windows.net
//
// For Azure Government or other sovereign clouds, the endpoint may differ.
type AzureConfig struct {
	// AccountName is the Azure Storage account name.
	// This is used to construct the full endpoint URL and for authentication.
	AccountName string

	// AccountKey is the primary or secondary access key for the Azure Storage account.
	// This key is used for Shared Key authentication with Azure Blob Storage.
	// The key should be kept secure and not exposed in logs or error messages.
	AccountKey string

	// Endpoint is the base URL for the Azure Blob Storage service.
	// For standard Azure Storage, this is typically:
	// https://{account-name}.blob.core.windows.net
	// For Azure Government or other sovereign clouds, the endpoint may differ.
	Endpoint string
}

// BucketProvider provides access to different types of blob storage buckets.
//
// The BucketProvider follows the factory pattern, creating different bucket
// instances based on the requested type. This allows for clear separation
// between public and private storage spaces while sharing common configuration.
//
// BucketProvider instances are safe for concurrent use.
type BucketProvider struct {
	// config holds the storage configuration used to create bucket instances.
	// The configuration is immutable after creation.
	config *Config
}

// New creates a new BucketProvider with the specified configuration.
//
// The function validates the configuration and prepares the provider for
// creating bucket instances. If the configuration is invalid, an error
// is returned with details about the validation failure.
//
// The returned BucketProvider is safe for concurrent use and can be used
// to create multiple bucket instances.
//
// Example:
//
//	config := &blob.Config{
//		Azure: blob.AzureConfig{
//			AccountName: "myaccount",
//			AccountKey:  "mykey",
//			Endpoint:    "https://myaccount.blob.core.windows.net",
//		},
//	}
//
//	provider, err := blob.New(config)
//	if err != nil {
//		log.Fatal(err)
//	}
func New(config *Config) (*BucketProvider, error) {
	if config == nil {
		return nil, fmt.Errorf("blob: config cannot be nil")
	}

	if config.Azure.AccountName == "" {
		return nil, fmt.Errorf("blob: Azure account name is required")
	}

	if config.Azure.AccountKey == "" {
		return nil, fmt.Errorf("blob: Azure account key is required")
	}

	if config.Azure.Endpoint == "" {
		return nil, fmt.Errorf("blob: Azure endpoint is required")
	}

	return &BucketProvider{config: config}, nil
}

// Public returns a bucket for public blob storage.
//
// The public bucket provides shared storage accessible to all users of the
// application. Files stored in the public bucket are typically accessible
// via public URLs and may be cached by CDNs or other services.
//
// The public bucket uses blob-level access control, allowing individual
// files to have different access permissions while maintaining the overall
// public nature of the bucket.
//
// Returns an error if the bucket cannot be created due to configuration
// issues or connectivity problems with the underlying storage service.
//
// Example:
//
//	publicBucket, err := provider.Public()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Upload a public file
//	file, _ := os.Open("public-image.jpg")
//	defer file.Close()
//	err = publicBucket.Upload(ctx, "images/logo.jpg", file, nil)
func (p *BucketProvider) Public() (*Bucket, error) {
	azConfig := &azurestore.AzConfig{
		AccountName:         p.config.Azure.AccountName,
		AccountKey:          p.config.Azure.AccountKey,
		Endpoint:            p.config.Azure.Endpoint,
		ContainerName:       PublicContainer,
		ContainerAccessType: blobAccessType,
		BlobAccessTier:      hotAccessTier,
	}

	azService, err := azurestore.NewAzureService(azConfig)
	if err != nil {
		return nil, fmt.Errorf("blob: failed to create Azure service: %w", err)
	}

	store := azurestore.New(azService)

	return &Bucket{b: store}, nil
}

// Space returns a bucket for space-specific blob storage.
//
// The space bucket provides isolated storage for a specific workspace or
// project identified by the spaceID parameter. Each space has its own
// container with private access control, ensuring data isolation between
// different workspaces.
//
// The spaceID is used to construct the container name in the format
// "space-{spaceID}". This naming convention ensures unique container
// names and clear identification of the storage space.
//
// Space buckets use private access control, requiring authentication
// for all operations. This provides security and isolation for
// workspace-specific data.
//
// Returns an error if the bucket cannot be created due to configuration
// issues, invalid spaceID, or connectivity problems with the underlying
// storage service.
//
// The spaceID parameter must be a valid string that can be used in
// Azure container names. Invalid characters will result in an error.
//
// Example:
//
//	spaceBucket, err := provider.Space("workspace-123")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Upload a private file
//	file, _ := os.Open("private-document.pdf")
//	defer file.Close()
//	err = spaceBucket.Upload(ctx, "documents/report.pdf", file, nil)
func (p *BucketProvider) Space(spaceID string) (*Bucket, error) {
	if spaceID == "" {
		return nil, fmt.Errorf("blob: spaceID cannot be empty")
	}

	azConfig := &azurestore.AzConfig{
		AccountName:         p.config.Azure.AccountName,
		AccountKey:          p.config.Azure.AccountKey,
		Endpoint:            p.config.Azure.Endpoint,
		ContainerName:       fmt.Sprintf("space-%s", spaceID),
		ContainerAccessType: privateAccessType,
		BlobAccessTier:      hotAccessTier,
	}

	azService, err := azurestore.NewAzureService(azConfig)
	if err != nil {
		return nil, fmt.Errorf("blob: failed to create Azure service: %w", err)
	}

	store := azurestore.New(azService)

	return &Bucket{b: store}, nil
}
