// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

// Package fga provides functions for managing Fine-Grained Authorization (FGA) models.
// It includes methods to load FGA models from files or DSL and register them with the FGA backend.
package fga

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/kopexa-grc/common/errors"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/openfga/language/pkg/go/transformer"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateModelFromFile loads an FGA model from a file and registers it with the FGA backend.
//
// If forceCreate is false and the latest model matches the file content, the existing model's
// ID is returned. Otherwise, a new model is created.
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
	// Load model from file
	dsl, err := os.ReadFile(fn)
	if err != nil {
		return "", err
	}

	if forceCreate {
		return c.createModelFromDSL(ctx, dsl)
	}

	return c.CreateModelFromDSL(ctx, dsl)
}

// CreateModelFromDSL creates or reuses an FGA model from a DSL definition.
//
// It first checks if the latest model in the store is identical to the provided DSL.
// If so, it returns the existing model's ID without writing a new one. This prevents
// unnecessary model proliferation when multiple instances start concurrently.
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

	// Check if the latest model already matches the new definition
	if modelID, err := c.findMatchingModel(ctx, body); err == nil && modelID != "" {
		return modelID, nil
	}

	return c.CreateModel(ctx, body)
}

// createModelFromDSL unconditionally creates a new FGA model from a DSL definition,
// bypassing the duplicate check.
func (c *Client) createModelFromDSL(ctx context.Context, dsl []byte) (string, error) {
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

// findMatchingModel checks whether the latest authorization model in the store
// is structurally identical to the given model request. Returns the existing
// model's ID if it matches, or an empty string if no match is found.
func (c *Client) findMatchingModel(ctx context.Context, newModel client.ClientWriteAuthorizationModelRequest) (string, error) {
	options := client.ClientReadAuthorizationModelsOptions{}

	resp, err := c.client.ReadAuthorizationModels(ctx).Options(options).Execute()
	if err != nil {
		return "", err
	}

	models := resp.GetAuthorizationModels()
	if len(models) == 0 {
		return "", nil
	}

	latest := models[0]

	if modelsEqual(latest, newModel) {
		log.Info().Str("model_id", latest.GetId()).Msg("fga model unchanged, reusing existing")
		return latest.GetId(), nil
	}

	return "", nil
}

// modelsEqual compares the type definitions and conditions of an existing
// authorization model with a new write request to determine if they are equivalent.
// It marshals both to JSON and compares the bytes for a structural comparison.
func modelsEqual(existing openfga.AuthorizationModel, newModel client.ClientWriteAuthorizationModelRequest) bool {
	existingTypes, err := json.Marshal(existing.GetTypeDefinitions())
	if err != nil {
		return false
	}

	newTypes, err := json.Marshal(newModel.GetTypeDefinitions())
	if err != nil {
		return false
	}

	if !bytes.Equal(existingTypes, newTypes) {
		return false
	}

	// Also compare conditions if present
	existingCond, err := json.Marshal(existing.GetConditions())
	if err != nil {
		return false
	}

	newCond, err := json.Marshal(newModel.GetConditions())
	if err != nil {
		return false
	}

	return bytes.Equal(existingCond, newCond)
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

// ExportDSLToJSON converts DSL bytes to a WriteAuthorizationModelRequest.
// This is exported for testing purposes.
func ExportDSLToJSON(dsl []byte) (client.ClientWriteAuthorizationModelRequest, error) {
	dslJSON, err := dslToJSON(dsl)
	if err != nil {
		return client.ClientWriteAuthorizationModelRequest{}, err
	}

	var body client.ClientWriteAuthorizationModelRequest
	if err := json.Unmarshal(dslJSON, &body); err != nil {
		return client.ClientWriteAuthorizationModelRequest{}, err
	}

	return body, nil
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
