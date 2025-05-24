// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga_test

import (
	"testing"

	"github.com/kopexa-grc/common/fga"
	"github.com/kopexa-grc/common/fga/internal/fgamock"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClient_CreateStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockList := fgamock.NewMockSdkClientListStoresRequestInterface(ctrl)
	mockCreate := fgamock.NewMockSdkClientCreateStoreRequestInterface(ctrl)

	tests := []struct {
		name           string
		storeName      string
		mockListResp   *client.ClientListStoresResponse
		mockCreateResp *client.ClientCreateStoreResponse
		mockListErr    error
		mockCreateErr  error
		expected       string
		expectError    bool
	}{
		{
			name:      "should return existing store ID",
			storeName: "test-store",
			mockListResp: &client.ClientListStoresResponse{
				Stores: []openfga.Store{
					{
						Id:   "store-123",
						Name: "existing-store",
					},
				},
			},
			expected:    "store-123",
			expectError: false,
		},
		{
			name:      "should create new store when none exist",
			storeName: "test-store",
			mockListResp: &client.ClientListStoresResponse{
				Stores: []openfga.Store{},
			},
			mockCreateResp: &client.ClientCreateStoreResponse{
				Id: "new-store-456",
			},
			expected:    "new-store-456",
			expectError: false,
		},
		{
			name:        "should handle list stores error",
			storeName:   "test-store",
			mockListErr: assert.AnError,
			expected:    "",
			expectError: true,
		},
		{
			name:      "should handle create store error",
			storeName: "test-store",
			mockListResp: &client.ClientListStoresResponse{
				Stores: []openfga.Store{},
			},
			mockCreateErr: assert.AnError,
			expected:      "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := fga.NewMockFGAClient(mockSdk)

			// Setup list stores mock
			mockSdk.EXPECT().ListStores(gomock.Any()).Return(mockList).Times(1)
			mockList.EXPECT().Options(gomock.Any()).Return(mockList).Times(1)
			mockList.EXPECT().Execute().Return(tt.mockListResp, tt.mockListErr).Times(1)

			// Setup create store mock if needed
			if tt.mockListErr == nil && len(tt.mockListResp.GetStores()) == 0 {
				mockSdk.EXPECT().CreateStore(gomock.Any()).Return(mockCreate).Times(1)
				mockCreate.EXPECT().Body(gomock.Any()).Return(mockCreate).Times(1)
				mockCreate.EXPECT().Execute().Return(tt.mockCreateResp, tt.mockCreateErr).Times(1)
			}

			storeID, err := c.CreateStore(tt.storeName)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, storeID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, storeID)
			}
		})
	}
}
