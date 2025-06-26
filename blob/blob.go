// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob

import (
	"fmt"

	"github.com/kopexa-grc/common/blob/azurestore"
)

type Config struct {
	Azure AzureConfig
}

type AzureConfig struct {
	AccountName string
	AccountKey  string
	Endpoint    string
}

type BucketProvider struct {
	config *Config
}

// New opens the bucket
func New(config *Config) (*BucketProvider, error) {
	return &BucketProvider{config: config}, nil
}

func (p *BucketProvider) Public() (*Bucket, error) {
	azConfig := &azurestore.AzConfig{
		AccountName:         p.config.Azure.AccountName,
		AccountKey:          p.config.Azure.AccountKey,
		Endpoint:            p.config.Azure.Endpoint,
		ContainerName:       "public",
		ContainerAccessType: blobAccessType,
		BlobAccessTier:      hotAccessTier,
	}

	azService, err := azurestore.NewAzureService(azConfig)
	if err != nil {
		return nil, err
	}

	store := azurestore.New(azService)

	return &Bucket{b: store}, nil
}

func (p *BucketProvider) Space(spaceID string) (*Bucket, error) {
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
		return nil, err
	}

	store := azurestore.New(azService)

	return &Bucket{b: store}, nil
}
