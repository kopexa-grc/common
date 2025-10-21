// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga_test

import (
	"context"
	"os"
	"testing"

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

	// Simulate: no model exists yet
	mockSdk.EXPECT().ReadAuthorizationModels(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Options(gomock.Any()).Return(mockRead).Times(1)
	mockRead.EXPECT().Execute().Return(&client.ClientReadAuthorizationModelsResponse{}, nil).Times(1)

	mockSdk.EXPECT().WriteAuthorizationModel(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Body(gomock.Any()).Return(mockWrite).Times(1)
	mockWrite.EXPECT().Execute().Return(&client.ClientWriteAuthorizationModelResponse{
		AuthorizationModelId: "test-model-id",
	}, nil).Times(1)

	modelPath := "testdata/model.fga"
	modelID, err := c.CreateModelFromFile(context.Background(), modelPath, true)
	assert.NoError(t, err)
	assert.Equal(t, "test-model-id", modelID)
}

func TestCreateModelFromDSL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockWrite := fgamock.NewMockSdkClientWriteAuthorizationModelRequestInterface(ctrl)

	c := fga.NewMockFGAClient(mockSdk)

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
