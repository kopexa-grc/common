// Original Licenses under Apache-2.0 by the openlane https://github.com/theopenlane
// SPDX-License-Identifier: Apache-2.0

package fga_test

import (
	"testing"

	"github.com/kopexa-grc/common/fga"
	"github.com/kopexa-grc/common/fga/internal/fgamock"
	"github.com/openfga/go-sdk/client"
	"go.uber.org/mock/gomock"
)

func TestClient_Grant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockWrite := fgamock.NewMockSdkClientWriteRequestInterface(ctrl)

	tests := []struct {
		name      string
		userID    string
		userType  string
		relation  string
		objType   string
		objID     string
		wantErr   bool
		setupMock func()
	}{
		{
			name:     "should grant permission successfully",
			userID:   "123",
			userType: "user",
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

			err := c.Grant().
				User(tt.userID).
				As(tt.userType).
				Relation(tt.relation).
				To(tt.objType, tt.objID).
				Apply(t.Context())

			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Grant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
