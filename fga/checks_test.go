package fga_test

import (
	"testing"

	"github.com/kopexa-grc/kopexa/pkg/fga"
	"github.com/kopexa-grc/kopexa/pkg/fga/internal/fgamock"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"go.uber.org/mock/gomock"
)

func TestClient_checkTuple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSdk := fgamock.NewMockSdkClient(ctrl)
	mockCheck := fgamock.NewMockSdkClientCheckRequestInterface(ctrl)

	type fields struct {
		user       string
		relation   string
		objectType string
		objectID   string
	}

	testsCases := []struct {
		name    string
		fields  fields
		want    bool
		wantErr bool
	}{
		{
			name: "should return true if the tuple exists",
			fields: fields{
				user:       "123",
				relation:   "member",
				objectType: "organization",
				objectID:   "kopexa",
			},
		},
	}
	for _, tt := range testsCases {
		t.Run(tt.name, func(t *testing.T) {

			c := fga.NewMockFGAClient(mockSdk)

			mockSdk.EXPECT().Check(gomock.Any()).Return(mockCheck).Times(1)
			mockCheck.EXPECT().Body(gomock.Any()).Return(mockCheck).Times(1)
			mockCheck.EXPECT().Execute().Return(&client.ClientCheckResponse{
				CheckResponse: openfga.CheckResponse{
					Allowed: &tt.want,
				},
			}, nil).Times(1)

			got, err := c.Has().User(tt.fields.user).Capability(tt.fields.relation).In(tt.fields.objectType, tt.fields.objectID).Check(t.Context())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.checkTuple() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.checkTuple() = %v, want %v", got, tt.want)
			}
		})
	}
}
