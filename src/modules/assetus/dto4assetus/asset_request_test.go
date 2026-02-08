package dto4assetus

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

func TestAssetRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     AssetRequest
		wantErr bool
	}{
		{"valid", AssetRequest{
			AssetID: "asset1",
			SpaceRequest: dto4spaceus.SpaceRequest{
				SpaceID: coretypes.SpaceID("space1"),
			},
		}, false},
		{"missing_asset_id", AssetRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{
				SpaceID: coretypes.SpaceID("space1"),
			},
		}, true},
		{"missing_space_id", AssetRequest{
			AssetID: "asset1",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AssetRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
