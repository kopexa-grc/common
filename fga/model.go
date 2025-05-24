// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

// Package fga provides functions for managing Fine-Grained Authorization (FGA) models.
// It includes methods to load FGA models from files or DSL and register them with the FGA backend.
package fga

import (
	"context"
	"encoding/json"
	"os"

	"github.com/kopexa-grc/common/errors"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/language/pkg/go/transformer"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateModelFromFile loads an FGA model from a file and registers it with the FGA backend.
//
// If forceCreate is false and a model already exists, the existing model's ID is returned.
// Otherwise, the model is loaded from the file and created anew.
//
// Parameters:
//   - ctx: Context for the request
//   - fn: Path to the model file (FGA DSL)
//   - forceCreate: true to always create a new model
//
// Returns:
//   - string: The model ID
//   - error: Any error that occurred
func (c *Client) CreateModelFromFile(ctx context.Context, fn string, forceCreate bool) (string, error) {
	options := client.ClientReadAuthorizationModelsOptions{}

	models, err := c.client.ReadAuthorizationModels(context.Background()).Options(options).Execute()
	if err != nil {
		return "", err
	}

	// Only create a new model if one does not exist and creation is not forced
	if !forceCreate {
		if len(models.AuthorizationModels) > 0 {
			modelID := models.GetAuthorizationModels()[0].Id
			log.Info().Str("model_id", modelID).Msg("fga model exists")

			return modelID, nil
		}
	}

	// Load model from file
	dsl, err := os.ReadFile(fn)
	if err != nil {
		return "", err
	}

	return c.CreateModelFromDSL(ctx, dsl)
}

// CreateModelFromDSL creates a new FGA model from a DSL definition ([]byte or string).
//
// Parameters:
//   - ctx: Context for the request
//   - dsl: FGA DSL as []byte
//
// Returns:
//   - string: The model ID
//   - error: Any error that occurred
func (c *Client) CreateModelFromDSL(ctx context.Context, dsl []byte) (string, error) {
	// Convert DSL to JSON
	dslJSON, err := dslToJSON(dsl)
	if err != nil {
		return "", err
	}

	var body client.ClientWriteAuthorizationModelRequest
	if err := json.Unmarshal(dslJSON, &body); err != nil {
		return "", err
	}

	return c.CreateModel(ctx, body)
}

// CreateModel registers an FGA model with the backend and returns the model ID.
//
// Parameters:
//   - ctx: Context for the request
//   - model: The model as WriteAuthorizationModelRequest
//
// Returns:
//   - string: The model ID
//   - error: Any error that occurred
func (c *Client) CreateModel(ctx context.Context, model client.ClientWriteAuthorizationModelRequest) (string, error) {
	resp, err := c.client.WriteAuthorizationModel(ctx).Body(model).Execute()
	if err != nil {
		return "", err
	}

	modelID := resp.GetAuthorizationModelId()

	log.Info().Str("model_id", modelID).Msg("fga model created")

	return modelID, nil
}

// dslToJSON converts an FGA model from DSL notation to JSON.
//
// Parameters:
//   - dslString: FGA DSL as []byte
//
// Returns:
//   - []byte: JSON representation
//   - error: Any error that occurred
func dslToJSON(dslString []byte) ([]byte, error) {
	parsedAuthModel, err := transformer.TransformDSLToProto(string(dslString))
	if err != nil {
		return []byte{}, errors.Wrap(err, ErrFailedToTransformModel.Error())
	}

	return protojson.Marshal(parsedAuthModel)
}
