// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"testing"

	"github.com/kopexa-grc/common/fga/internal/fgamock"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/go-sdk/credentials"

	"github.com/stretchr/testify/assert"
)

func NewMockFGAClient(c *fgamock.MockSdkClient) *Client {
	return &Client{
		client: c,
	}
}

const mockStoreID = "01JV5FY6B75PMFSK86MV6EX3Y9"

func TestNewClient_Success(t *testing.T) {
	client, err := NewClient("https://api.openfga.example")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.openfga.example", client.config.ApiUrl)
	assert.True(t, client.IgnoreDuplicateKeyError)
}

func TestNewClient_WithOptions(t *testing.T) {
	client, err := NewClient(
		"https://api.openfga.example",
		WithStoreID(mockStoreID),
		WithIgnoreDuplicateKeyError(false),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.openfga.example", client.config.ApiUrl)
	assert.Equal(t, mockStoreID, client.config.StoreId)
	assert.False(t, client.IgnoreDuplicateKeyError)
}

func TestNewClient_EmptyHost(t *testing.T) {
	client, err := NewClient("")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestWithStoreID(t *testing.T) {
	c := &Client{config: &client.ClientConfiguration{}}
	opt := WithStoreID(mockStoreID)
	opt(c)
	assert.Equal(t, mockStoreID, c.config.StoreId)
}

func TestWithIgnoreDuplicateKeyError(t *testing.T) {
	c := &Client{}
	opt := WithIgnoreDuplicateKeyError(true)
	opt(c)
	assert.True(t, c.IgnoreDuplicateKeyError)
}

func TestWithToken(t *testing.T) {
	c := &Client{config: &client.ClientConfiguration{}}
	opt := WithToken("test-token")
	opt(c)
	assert.Equal(t, credentials.CredentialsMethodApiToken, c.config.Credentials.Method)
	assert.Equal(t, "test-token", c.config.Credentials.Config.ApiToken)
}
