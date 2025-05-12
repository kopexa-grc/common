// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package azurestore_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/kopexa-grc/common/blob/azurestore"
	"github.com/kopexa-grc/common/blob/driver"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

const mockContainer = "kopexa"
const mockID = "123"

var mockInfo = map[string]interface{}{
	"id":   mockID,
	"name": "test",
	"size": 100,
	"type": "file",
	"url":  "https://example.com/test",
}

func TestGetSignedUploadURL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	assert := assert.New(t)

	ctx := context.Background()
	mockKey := "avatar123.png"
	expectedURL := "https://storage.example.com/signed-url"

	mockService := NewMockAzService(mockCtrl)
	mockBlob := NewMockAzBlob(mockCtrl)

	gomock.InOrder(
		mockService.EXPECT().
			NewBlob(ctx, mockKey).
			Return(mockBlob, nil).
			Times(1),
		mockBlob.EXPECT().
			SignedURL(ctx, &driver.SignedURLOptions{
				Expiry:      time.Minute * 15,
				Method:      http.MethodPut,
				ContentType: "image/png",
			}).
			Return(expectedURL, nil).
			Times(1),
	)

	store := azurestore.New(mockService)
	store.Container = "avatars"

	url, err := store.GetSignedUploadURL(ctx, mockKey, time.Minute*15, 5*1024*1024, "image/png")
	assert.NoError(err)
	assert.Equal(expectedURL, url)
}

func TestGetSignedDownloadURL(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	assert := assert.New(t)

	ctx := context.Background()
	mockKey := "avatar123.png"
	expectedURL := "https://storage.example.com/signed-url"

	mockService := NewMockAzService(mockCtrl)
	mockBlob := NewMockAzBlob(mockCtrl)

	gomock.InOrder(
		mockService.EXPECT().
			NewBlob(ctx, mockKey).
			Return(mockBlob, nil).
			Times(1),
		mockBlob.EXPECT().
			SignedURL(ctx, &driver.SignedURLOptions{
				Expiry: time.Minute * 15,
				Method: http.MethodGet,
			}).
			Return(expectedURL, nil).
			Times(1),
	)

	store := azurestore.New(mockService)
	store.Container = "avatars"

	url, err := store.GetSignedDownloadURL(ctx, mockKey, time.Minute*15)
	assert.NoError(err)
	assert.Equal(expectedURL, url)
}

func TestDelete(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	assert := assert.New(t)

	ctx, cancel := context.WithCancel(t.Context())

	service := NewMockAzService(mockCtrl)
	store := azurestore.New(service)
	store.Container = mockContainer

	blockBlob := NewMockAzBlob(mockCtrl)
	assert.NotNil(blockBlob)

	gomock.InOrder(
		service.EXPECT().NewBlob(ctx, mockID+".info").Return(blockBlob, nil).Times(1),
		blockBlob.EXPECT().Delete(ctx).Return(nil).Times(1),
	)

	store.DeleteObject(ctx, mockID+".info")

	cancel()
}
