// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package blob

import (
	"context"
	"testing"

	"github.com/kopexa-grc/common/blob/driver"
	"github.com/stretchr/testify/assert"
)

// mockBucket implements driver.Bucket for testing
type mockBucket struct {
	deleteFunc    func(ctx context.Context, key string) error
	signedURLFunc func(ctx context.Context, key string, opts *driver.SignedURLOptions) (string, error)
	copyFunc      func(ctx context.Context, dstKey, srcKey string, opts *driver.CopyOptions) error
}

func (m *mockBucket) Delete(ctx context.Context, key string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, key)
	}
	return nil
}

func (m *mockBucket) SignedURL(ctx context.Context, key string, opts *driver.SignedURLOptions) (string, error) {
	if m.signedURLFunc != nil {
		return m.signedURLFunc(ctx, key, opts)
	}
	return "https://mock-url.com/" + key, nil
}

func (m *mockBucket) Copy(ctx context.Context, dstKey, srcKey string, opts *driver.CopyOptions) error {
	if m.copyFunc != nil {
		return m.copyFunc(ctx, dstKey, srcKey, opts)
	}
	return nil
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Azure: AzureConfig{
					AccountName: "test-account",
					AccountKey:  "dGVzdC1rZXk=", // base64 encoded "test-key"
					Endpoint:    "https://test.blob.core.windows.net",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, provider)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)
			}
		})
	}
}
