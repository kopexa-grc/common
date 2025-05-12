// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/kopexa-grc/common/blob"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBucket_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDriver := NewMockBucket(ctrl)
	bucket := blob.NewBucketForTest(mockDriver)

	tests := []struct {
		name    string
		key     string
		setup   func()
		wantErr bool
	}{
		{
			name:    "invalid UTF-8 key",
			key:     string([]byte{0xFF, 0xFE, 0xFD}),
			setup:   func() {},
			wantErr: true,
		},
		{
			name:    "empty key",
			key:     "",
			setup:   func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := bucket.Delete(context.Background(), tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBucket_SignedURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDriver := NewMockBucket(ctrl)
	bucket := blob.NewBucketForTest(mockDriver)

	tests := []struct {
		name    string
		key     string
		opts    *blob.SignedURLOptions
		setup   func()
		wantErr bool
	}{
		{
			name: "valid GET URL",
			key:  "test-key",
			opts: &blob.SignedURLOptions{
				Method: http.MethodGet,
				Expiry: time.Hour,
			},
			setup: func() {
				mockDriver.EXPECT().
					SignedURL(gomock.Any(), "test-key", gomock.Any()).
					Return("https://test-url.com/test-key", nil)
			},
			wantErr: false,
		},
		{
			name: "valid PUT URL",
			key:  "test-key",
			opts: &blob.SignedURLOptions{
				Method:      http.MethodPut,
				Expiry:      time.Hour,
				ContentType: "application/json",
			},
			setup: func() {
				mockDriver.EXPECT().
					SignedURL(gomock.Any(), "test-key", gomock.Any()).
					Return("https://test-url.com/test-key", nil)
			},
			wantErr: false,
		},
		{
			name: "invalid method",
			key:  "test-key",
			opts: &blob.SignedURLOptions{
				Method: "INVALID",
				Expiry: time.Hour,
			},
			setup:   func() {},
			wantErr: true,
		},
		{
			name: "invalid UTF-8 key",
			key:  string([]byte{0xFF, 0xFE, 0xFD}),
			opts: &blob.SignedURLOptions{
				Method: http.MethodGet,
				Expiry: time.Hour,
			},
			setup:   func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			url, err := bucket.SignedURL(context.Background(), tt.key, tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, url)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, url)
			}
		})
	}
}

func TestBucket_Copy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDriver := NewMockBucket(ctrl)
	bucket := blob.NewBucketForTest(mockDriver)

	tests := []struct {
		name    string
		dstKey  string
		srcKey  string
		opts    *blob.CopyOptions
		setup   func()
		wantErr bool
	}{
		{
			name:   "valid copy",
			dstKey: "dst-key",
			srcKey: "src-key",
			opts:   &blob.CopyOptions{},
			setup: func() {
				mockDriver.EXPECT().
					Copy(gomock.Any(), "dst-key", "src-key", gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "invalid UTF-8 dst key",
			dstKey:  string([]byte{0xFF, 0xFE, 0xFD}),
			srcKey:  "src-key",
			opts:    &blob.CopyOptions{},
			setup:   func() {},
			wantErr: true,
		},
		{
			name:    "invalid UTF-8 src key",
			dstKey:  "dst-key",
			srcKey:  string([]byte{0xFF, 0xFE, 0xFD}),
			opts:    &blob.CopyOptions{},
			setup:   func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			err := bucket.Copy(context.Background(), tt.dstKey, tt.srcKey, tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
