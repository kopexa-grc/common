package fga_test

import (
	"context"
	"errors"
	"testing"

	"github.com/kopexa-grc/kopexa/pkg/fga"
	"github.com/kopexa-grc/kopexa/pkg/fga/internal/fgamock"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestClient_ListTuples(t *testing.T) {
	tests := []struct {
		name        string
		req         fga.ListTuplesRequest
		mockResp    *client.ClientReadResponse
		mockErr     error
		expected    *fga.ListTuplesResponse
		expectError bool
	}{
		{
			name: "successful list tuples",
			req: fga.ListTuplesRequest{
				Subject: fga.Entity{
					Kind:       "user",
					Identifier: "123",
				},
				Relation: "member",
				Object: fga.Entity{
					Kind:       "organization",
					Identifier: "456",
				},
			},
			mockResp: &client.ClientReadResponse{
				Tuples: []openfga.Tuple{
					{
						Key: openfga.TupleKey{
							User:     "user:123",
							Relation: "member",
							Object:   "organization:456",
						},
					},
				},
			},
			expected: &fga.ListTuplesResponse{
				Tuples: []fga.TupleKey{
					{
						Subject: fga.Entity{
							Kind:       "user",
							Identifier: "123",
						},
						Relation: "member",
						Object: fga.Entity{
							Kind:       "organization",
							Identifier: "456",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:        "empty request",
			req:         fga.ListTuplesRequest{},
			expected:    nil,
			expectError: true,
		},
		{
			name: "client error",
			req: fga.ListTuplesRequest{
				Subject: fga.Entity{
					Kind:       "user",
					Identifier: "123",
				},
			},
			mockErr:     errors.New("client error"),
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSdk := fgamock.NewMockSdkClient(ctrl)
			mockRead := fgamock.NewMockSdkClientReadRequestInterface(ctrl)
			c := fga.NewMockFGAClient(mockSdk)

			if tt.mockErr == nil && tt.req.Subject.Identifier != "" {
				mockSdk.EXPECT().Read(gomock.Any()).Return(mockRead).Times(1)
				mockRead.EXPECT().Body(gomock.Any()).Return(mockRead).Times(1)
				mockRead.EXPECT().Execute().Return(tt.mockResp, nil).Times(1)
			} else if tt.mockErr != nil {
				mockSdk.EXPECT().Read(gomock.Any()).Return(mockRead).Times(1)
				mockRead.EXPECT().Body(gomock.Any()).Return(mockRead).Times(1)
				mockRead.EXPECT().Execute().Return(nil, tt.mockErr).Times(1)
			}

			resp, err := c.ListTuples(context.Background(), tt.req)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp)
			}
		})
	}
}

func TestClient_WriteTupleKeys(t *testing.T) {
	tests := []struct {
		name        string
		writes      []fga.TupleKey
		deletes     []fga.TupleKey
		mockResp    *client.ClientWriteResponse
		mockErr     error
		expected    *client.ClientWriteResponse
		expectError bool
	}{
		{
			name: "successful write",
			writes: []fga.TupleKey{
				{
					Subject: fga.Entity{
						Kind:       "user",
						Identifier: "123",
					},
					Relation: "member",
					Object: fga.Entity{
						Kind:       "organization",
						Identifier: "456",
					},
				},
			},
			mockResp: &client.ClientWriteResponse{
				Writes: []client.ClientWriteRequestWriteResponse{
					{
						TupleKey: client.ClientTupleKey{
							User:     "user:123",
							Relation: "member",
							Object:   "organization:456",
						},
					},
				},
			},
			expected: &client.ClientWriteResponse{
				Writes: []client.ClientWriteRequestWriteResponse{
					{
						TupleKey: client.ClientTupleKey{
							User:     "user:123",
							Relation: "member",
							Object:   "organization:456",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "client error",
			writes: []fga.TupleKey{
				{
					Subject: fga.Entity{
						Kind:       "user",
						Identifier: "123",
					},
					Relation: "member",
					Object: fga.Entity{
						Kind:       "organization",
						Identifier: "456",
					},
				},
			},
			mockErr:     errors.New("client error"),
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSdk := fgamock.NewMockSdkClient(ctrl)
			mockWrite := fgamock.NewMockSdkClientWriteRequestInterface(ctrl)
			c := fga.NewMockFGAClient(mockSdk)

			if tt.mockErr == nil {
				mockSdk.EXPECT().Write(gomock.Any()).Return(mockWrite).Times(1)
				mockWrite.EXPECT().Body(gomock.Any()).Return(mockWrite).Times(1)
				mockWrite.EXPECT().Options(gomock.Any()).Return(mockWrite).Times(1)
				mockWrite.EXPECT().Execute().Return(tt.mockResp, nil).Times(1)
			} else {
				mockSdk.EXPECT().Write(gomock.Any()).Return(mockWrite).Times(1)
				mockWrite.EXPECT().Body(gomock.Any()).Return(mockWrite).Times(1)
				mockWrite.EXPECT().Options(gomock.Any()).Return(mockWrite).Times(1)
				mockWrite.EXPECT().Execute().Return(nil, tt.mockErr).Times(1)
			}

			resp, err := c.WriteTupleKeys(context.Background(), tt.writes, tt.deletes)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, resp)
			}
		})
	}
}
