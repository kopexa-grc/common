// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga

import (
	"context"

	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/rs/zerolog/log"
)

// CreateStore creates a new fine-grained authorization store or returns an existing one.
// This function implements a "create if not exists" pattern for FGA stores.
//
// The function:
// 1. Checks for existing stores
// 2. Returns the first store's ID if any exist
// 3. Creates a new store with the given name if no stores exist
//
// Example:
//
//	storeID, err := client.CreateStore("my-auth-store")
//	if err != nil {
//	    // Handle error
//	}
//
// Parameters:
//   - storeName: The name to use for the new store if one needs to be created
//
// Returns:
//   - string: The ID of the store (either existing or newly created)
//   - error: If the store creation or listing fails
func (c *Client) CreateStore(storeName string) (string, error) {
	options := client.ClientListStoresOptions{
		ContinuationToken: openfga.PtrString(""),
	}

	stores, err := c.client.ListStores(context.Background()).Options(options).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to list FGA stores")

		return "", err
	}

	// Only create a new store if one does not exist
	if len(stores.GetStores()) > 0 {
		storeID := stores.GetStores()[0].Id
		log.Info().
			Str("store_id", storeID).
			Msg("using existing FGA store")

		return storeID, nil
	}

	// Create new store
	storeReq := c.client.CreateStore(context.Background())

	resp, err := storeReq.Body(client.ClientCreateStoreRequest{
		Name: storeName,
	}).Execute()
	if err != nil {
		log.Error().
			Err(err).
			Str("store_name", storeName).
			Msg("failed to create FGA store")

		return "", err
	}

	storeID := resp.GetId()

	log.Info().
		Str("store_id", storeID).
		Str("store_name", storeName).
		Msg("created new FGA store")

	return storeID, nil
}
