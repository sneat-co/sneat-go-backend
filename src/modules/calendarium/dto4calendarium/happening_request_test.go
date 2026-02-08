package dto4calendarium

import (
	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"testing"
)

func TestHappeningRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     HappeningRequest
		wantErr bool
	}{
		{"valid", HappeningRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
			HappeningID:  "h1",
		}, false},
		{"valid_with_type", HappeningRequest{
			SpaceRequest:  dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
			HappeningID:   "h1",
			HappeningType: dbo4calendarium.HappeningTypeSingle,
		}, false},
		{"missing_id", HappeningRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
		}, true},
		{"invalid_type", HappeningRequest{
			SpaceRequest:  dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("s1")},
			HappeningID:   "h1",
			HappeningType: "invalid",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("HappeningRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
