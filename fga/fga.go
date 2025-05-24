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
	"context"

	"github.com/kopexa-grc/common/errors"
	"github.com/openfga/go-sdk/client"
)

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

// CreateClientWithStore creates a new FGA client with a store and model configuration.
// It handles the complete setup process including:
// - Creating or using an existing store
// - Setting up authentication credentials
// - Creating or using an existing authorization model
//
// The function supports two ways of providing the model:
// 1. Using ModelData: Directly providing the model definition as []byte
// 2. Using ModelFile: Loading the model from a file
//
// Example:
//
//	config := fga.Config{
//	    HostURL: "https://api.openfga.example",
//	    ModelData: []byte("model\n  schema 1.1\n\ntype user\n\ntype document\n  relations\n    define viewer: [user]"),
//	    StoreName: "my-store",
//	}
//	client, err := fga.CreateClientWithStore(ctx, config)
//
// Parameters:
//   - ctx: Context for the request
//   - c: Configuration containing store and model settings
//
// Returns:
//   - *Client: A configured FGA client
//   - error: If the setup process fails
func CreateClientWithStore(ctx context.Context, c Config) (*Client, error) {
	opts := []Option{
		WithIgnoreDuplicateKeyError(c.IgnoreDuplicateKeyError),
	}

	// set credentials if provided
	if c.Credentials.APIToken != "" {
		opts = append(opts, WithAPITokenCredentials(c.Credentials.APIToken))
	} else if c.Credentials.ClientID != "" && c.Credentials.ClientSecret != "" {
		opts = append(opts, WithClientCredentials(
			c.Credentials.ClientID,
			c.Credentials.ClientSecret,
			c.Credentials.Audience,
			c.Credentials.Issuer,
			c.Credentials.Scopes,
		))
	}

	// create store if an ID was not configured
	if c.StoreID == "" {
		// Create new store
		fgaClient, err := NewClient(
			c.HostURL,
			opts...,
		)
		if err != nil {
			return nil, err
		}

		c.StoreID, err = fgaClient.CreateStore(c.StoreName)
		if err != nil {
			return nil, err
		}
	}

	// add store ID to the options
	opts = append(opts, WithStoreID(c.StoreID))

	// create model if ID was not configured
	if c.ModelID == "" {
		// create fga client with store ID
		fgaClient, err := NewClient(
			c.HostURL,
			opts...,
		)
		if err != nil {
			return nil, err
		}

		var modelID string
		if c.ModelData != nil {
			// Create model from provided data
			modelID, err = fgaClient.CreateModelFromDSL(ctx, c.ModelData)
		} else {
			// Create model from file if no data provided
			modelID, err = fgaClient.CreateModelFromFile(ctx, c.ModelFile, c.CreateNewModel)
		}

		if err != nil {
			return nil, err
		}

		// Set ModelID in the config
		c.ModelID = modelID
	}

	// add model ID to the options
	opts = append(opts,
		WithAuthorizationModelID(c.ModelID),
	)

	// create fga client with store ID
	return NewClient(
		c.HostURL,
		opts...,
	)
}

// Healthcheck returns a function that checks if the FGA service is accessible.
// The returned function can be used as a health check in service monitoring.
// It verifies the connection by attempting to read the authorization model.
//
// Example:
//
//	healthcheck := fga.Healthcheck(client)
//	err := healthcheck(ctx)
//	if err != nil {
//	    // Service is not healthy
//	}
//
// Parameters:
//   - c: The FGA client to check
//
// Returns:
//   - func(ctx context.Context) error: A function that performs the health check
func Healthcheck(c Client) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		opts := client.ClientReadAuthorizationModelOptions{
			AuthorizationModelId: &c.config.AuthorizationModelId,
		}

		_, err := c.client.ReadAuthorizationModel(ctx).Options(opts).Execute()
		if err != nil {
			return err
		}

		return nil
	}
}
