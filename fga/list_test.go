// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga_test

import (
	"context"
	"testing"

	"github.com/kopexa-grc/common/fga"
	"github.com/kopexa-grc/common/fga/internal/fgamock"
	"github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClient_ListObjectIDsWithAccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockList := fgamock.NewMockSdkClientListObjectsRequestInterface(ctrl)

	tests := []struct {
		name        string
		req         fga.ListRequest
		mockResp    *client.ClientListObjectsResponse
		mockErr     error
		expected    []string
		expectError bool
	}{
		{
			name: "should return list of object IDs successfully",
			req: fga.ListRequest{
				SubjectID:   "user123",
				SubjectType: "user",
				ObjectType:  "document",
				Relation:    "viewer",
			},
			mockResp: &client.ClientListObjectsResponse{
				Objects: []string{
					"document:doc1",
					"document:doc2",
					"document:doc3",
				},
			},
			expected:    []string{"doc1", "doc2", "doc3"},
			expectError: false,
		},
		{
			name: "should handle empty response",
			req: fga.ListRequest{
				SubjectID:   "user123",
				SubjectType: "user",
				ObjectType:  "document",
				Relation:    "viewer",
			},
			mockResp: &client.ClientListObjectsResponse{
				Objects: []string{},
			},
			expected:    []string{},
			expectError: false,
		},
		{
			name: "should handle FGA service error",
			req: fga.ListRequest{
				SubjectID:   "user123",
				SubjectType: "user",
				ObjectType:  "document",
				Relation:    "viewer",
			},
			mockErr:     assert.AnError,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := fga.NewMockFGAClient(mockSdk)

			if tt.mockErr == nil {
				mockSdk.EXPECT().ListObjects(gomock.Any()).Return(mockList).Times(1)
				mockList.EXPECT().Body(gomock.Any()).Return(mockList).Times(1)
				mockList.EXPECT().Execute().Return(tt.mockResp, nil).Times(1)
			} else {
				mockSdk.EXPECT().ListObjects(gomock.Any()).Return(mockList).Times(1)
				mockList.EXPECT().Body(gomock.Any()).Return(mockList).Times(1)
				mockList.EXPECT().Execute().Return(nil, tt.mockErr).Times(1)
			}

			objectIDs, err := c.ListObjectIDsWithAccess(context.Background(), tt.req)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, objectIDs)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, objectIDs)
			}
		})
	}
}
