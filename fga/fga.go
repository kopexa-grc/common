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
)

// Option is a function that configures a Client.
// Options are used to customize the behavior of the FGA client.
type Option func(*Client)

// Options contains configuration options for the FGA client.
// These options are used to configure the client's behavior.
type Options struct {
	storeID string
}

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
