// Copyright (c) Kopexa GmbH
// SPDX-License-Identifier: BUSL-1.1

package fga_test

import (
	"testing"

	"github.com/kopexa-grc/common/fga"
	"github.com/kopexa-grc/common/fga/internal/fgamock"
	"github.com/openfga/go-sdk/client"
	"go.uber.org/mock/gomock"
)

func TestClient_Revoke(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockWrite := fgamock.NewMockSdkClientWriteRequestInterface(ctrl)

	tests := []struct {
		name      string
		userID    string
		relation  string
		objType   string
		objID     string
		wantErr   bool
		setupMock func()
	}{
		{
			name:     "should revoke permission successfully",
			userID:   "123",
			relation: "member",
			objType:  "organization",
			objID:    "kopexa",
			wantErr:  false,
			setupMock: func() {
				mockSdk.EXPECT().Write(gomock.Any()).Return(mockWrite)
				mockWrite.EXPECT().Body(gomock.Any()).Return(mockWrite)
				mockWrite.EXPECT().Options(gomock.Any()).Return(mockWrite)
				mockWrite.EXPECT().Execute().Return(&client.ClientWriteResponse{
					Writes:  []client.ClientWriteRequestWriteResponse{},
					Deletes: []client.ClientWriteRequestDeleteResponse{},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := fga.NewMockFGAClient(mockSdk)

			tt.setupMock()

			err := c.Revoke().
				User(tt.userID).
				Relation(tt.relation).
				From(tt.objType, tt.objID).
				Apply(t.Context())

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Revoke() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
