// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga_test

import (
	"context"
	"os"
	"testing"

	openfga "github.com/openfga/go-sdk"

	"github.com/kopexa-grc/common/fga"
	"github.com/kopexa-grc/common/fga/internal/fgamock"
	"github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateModelFromFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockRead := fgamock.NewMockSdkClientReadAuthorizationModelsRequestInterface(ctrl)
	mockWrite := fgamock.NewMockSdkClientWriteAuthorizationModelRequestInterface(ctrl)

	c := fga.NewMockFGAClient(mockSdk)

	// CreateModelFromFile with forceCreate=true should bypass the check
	// and go directly to createModelFromDSL (no ReadAuthorizationModels call)
	mockSdk.EXPECT().WriteAuthorizationModel(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Body(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Execute().Return(&client.ClientWriteAuthorizationModelResponse{
		AuthorizationModelId: "test-model-id",
	}, nil).Times(1)

	modelPath := "testdata/model.fga"
	modelID, err := c.CreateModelFromFile(context.Background(), modelPath, true)
	assert.NoError(t, err)
	assert.Equal(t, "test-model-id", modelID)

	// Now test with forceCreate=false — should check for existing models first.
	// Since no models exist, it should write a new one.
	mockSdk.EXPECT().ReadAuthorizationModels(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Options(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Execute().Return(&client.ClientReadAuthorizationModelsResponse{}, nil).Times(1)

	mockSdk.EXPECT().WriteAuthorizationModel(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Body(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Execute().Return(&client.ClientWriteAuthorizationModelResponse{
		AuthorizationModelId: "new-model-id",
	}, nil).Times(1)

	modelID, err = c.CreateModelFromFile(context.Background(), modelPath, false)
	assert.NoError(t, err)
	assert.Equal(t, "new-model-id", modelID)
}

func TestCreateModelFromDSL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockRead := fgamock.NewMockSdkClientReadAuthorizationModelsRequestInterface(ctrl)
	mockWrite := fgamock.NewMockSdkClientWriteAuthorizationModelRequestInterface(ctrl)

	c := fga.NewMockFGAClient(mockSdk)

	// Simulate: no models exist yet → should write new model
	mockSdk.EXPECT().ReadAuthorizationModels(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Options(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Execute().Return(&client.ClientReadAuthorizationModelsResponse{}, nil).Times(1)

	mockSdk.EXPECT().WriteAuthorizationModel(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Body(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Execute().Return(&client.ClientWriteAuthorizationModelResponse{
		AuthorizationModelId: "test-model-id",
	}, nil).Times(1)

	dsl, err := os.ReadFile("testdata/model.fga")
	assert.NoError(t, err)

	modelID, err := c.CreateModelFromDSL(context.Background(), dsl)
	assert.NoError(t, err)
	assert.Equal(t, "test-model-id", modelID)
}

func TestCreateModelFromDSL_ReusesExistingModel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockRead := fgamock.NewMockSdkClientReadAuthorizationModelsRequestInterface(ctrl)

	c := fga.NewMockFGAClient(mockSdk)

	dsl, err := os.ReadFile("testdata/model.fga")
	assert.NoError(t, err)

	// Build the expected type definitions from the test DSL so we can return
	// them as the "existing model" from ReadAuthorizationModels.
	existingModel := buildModelFromDSL(t, dsl)

	// Simulate: existing model matches → should NOT call WriteAuthorizationModel
	mockSdk.EXPECT().ReadAuthorizationModels(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Options(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Execute().Return(&client.ClientReadAuthorizationModelsResponse{
		AuthorizationModels: []openfga.AuthorizationModel{existingModel},
	}, nil).Times(1)

	// No WriteAuthorizationModel call expected!

	modelID, err := c.CreateModelFromDSL(context.Background(), dsl)
	assert.NoError(t, err)
	assert.Equal(t, "existing-model-id", modelID)
}

func TestCreateModelFromDSL_WritesWhenModelDiffers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockRead := fgamock.NewMockSdkClientReadAuthorizationModelsRequestInterface(ctrl)
	mockWrite := fgamock.NewMockSdkClientWriteAuthorizationModelRequestInterface(ctrl)

	c := fga.NewMockFGAClient(mockSdk)

	dsl, err := os.ReadFile("testdata/model.fga")
	assert.NoError(t, err)

	// Return a model with different type definitions
	differentModel := openfga.AuthorizationModel{
		Id:            "old-model-id",
		SchemaVersion: "1.1",
		TypeDefinitions: []openfga.TypeDefinition{
			{Type: "user"},
			// Missing "document" type → different from DSL
		},
	}

	// Simulate: existing model differs → should write new model
	mockSdk.EXPECT().ReadAuthorizationModels(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Options(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Execute().Return(&client.ClientReadAuthorizationModelsResponse{
		AuthorizationModels: []openfga.AuthorizationModel{differentModel},
	}, nil).Times(1)

	mockSdk.EXPECT().WriteAuthorizationModel(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Body(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Execute().Return(&client.ClientWriteAuthorizationModelResponse{
		AuthorizationModelId: "new-model-id",
	}, nil).Times(1)

	modelID, err := c.CreateModelFromDSL(context.Background(), dsl)
	assert.NoError(t, err)
	assert.Equal(t, "new-model-id", modelID)
}

// buildModelFromDSL converts DSL bytes into an AuthorizationModel matching
// what OpenFGA would return from ReadAuthorizationModels.
func buildModelFromDSL(t *testing.T, dsl []byte) openfga.AuthorizationModel {
	t.Helper()

	// Use the same DSL→JSON→request path as the production code
	parsed, err := fga.ExportDSLToJSON(dsl)
	assert.NoError(t, err)

	return openfga.AuthorizationModel{
		Id:              "existing-model-id",
		SchemaVersion:   parsed.SchemaVersion,
		TypeDefinitions: parsed.TypeDefinitions,
		Conditions:      parsed.Conditions,
	}
}
