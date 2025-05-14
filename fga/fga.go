// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package fga provides a high-level client for interacting with the OpenFGA service.
// It implements a fluent interface for managing fine-grained authorization (FGA) in a type-safe manner.
//
// The package is designed to be:
//   - Type-safe: All operations are strongly typed to prevent runtime errors
//   - Fluent: Uses a builder pattern for clear and readable API calls
//   - Efficient: Minimizes allocations and network calls
//   - Reliable: Handles errors gracefully and provides detailed error information
//
// Example usage:
//
//	client, err := fga.NewClient("https://api.openfga.example")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check if a user has access
//	hasAccess, err := client.Has().
//	    User("user123").
//	    Capability("viewer").
//	    In("document", "doc123").
//	    Check(ctx)
//
//	// Grant access
//	err = client.Grant().
//	    User("user123").
//	    Relation("viewer").
//	    To("document", "doc123").
//	    Apply(ctx)
package fga

import (
	"github.com/kopexa-grc/common/errors"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/go-sdk/credentials"
)

// Option is a function that configures a Client.
// Options are used to customize the behavior of the FGA client.
type Option func(*Client)

// WithStoreID sets the store ID for the FGA client.
// This is required when working with multiple stores.
// The store ID is used to identify which authorization store to use.
func WithStoreID(storeID string) Option {
	return func(c *Client) {
		c.config.StoreId = storeID
	}
}

// WithIgnoreDuplicateKeyError configures whether duplicate key errors should be ignored.
// When set to true, attempts to write duplicate tuples will be silently ignored.
// This is useful in scenarios where idempotency is desired.
func WithIgnoreDuplicateKeyError(ignore bool) Option {
	return func(c *Client) {
		c.IgnoreDuplicateKeyError = ignore
	}
}

// Credentials represents the authentication credentials for the OpenFGA service.
// It is used to configure the client with the necessary authentication information.
type Credentials struct {
	// APIToken is the API token used to authenticate with the OpenFGA service.
	// This token is required for all API calls to the service.
	APIToken string `json:"api_token" koanf:"api_token" jsonschema:"description=The API token for the OpenFGA service"`
}

// WithToken configures the FGA client with an API token for authentication.
// The token is used to authenticate all requests to the OpenFGA service.
// This option is required for production use of the client.
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithToken("your-api-token"),
//	)
func WithToken(token string) Option {
	return func(c *Client) {
		c.config.Credentials = &credentials.Credentials{
			Method: credentials.CredentialsMethodApiToken,
			Config: &credentials.Config{
				ApiToken: token,
			},
		}
	}
}

// Client represents a connection to the OpenFGA service.
// It provides methods for checking and managing permissions in a type-safe manner.
// The client is safe for concurrent use by multiple goroutines.
type Client struct {
	client client.SdkClient
	config *client.ClientConfiguration

	// IgnoreDuplicateKeyError determines whether duplicate key errors should be ignored.
	// When true, attempts to write duplicate tuples will be silently ignored.
	IgnoreDuplicateKeyError bool
}

// NewClient creates a new FGA client with the given host and options.
// The host parameter is required and should be the URL of the OpenFGA service.
// Returns an error if the client cannot be created.
//
// Example:
//
//	client, err := fga.NewClient("https://api.openfga.example",
//	    fga.WithStoreID("store123"),
//	    fga.WithIgnoreDuplicateKeyError(true),
//	)
func NewClient(host string, opts ...Option) (*Client, error) {
	var err error

	if host == "" {
		return nil, errors.NewInvalidArgument("host is required")
	}

	c := &Client{
		config: &client.ClientConfiguration{
			ApiUrl: host,
		},
		IgnoreDuplicateKeyError: true,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.client, err = client.NewSdkClient(c.config)
	if err != nil {
		return nil, err
	}

	return c, nil
}
