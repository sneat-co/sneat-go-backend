package dto4calendarius

import (
	"testing"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/stretchr/testify/assert"
)

// validHappeningRequest returns a minimal valid HappeningRequest.
func validHappeningRequest() HappeningRequest {
	return HappeningRequest{
		SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("space1")},
		HappeningID:  "happening1",
	}
}

// validSlotTiming mirrors a valid HappeningSlotTiming for use in DTO tests.
func validSlotTiming() dbo4calendarius.HappeningSlotTiming {
	return dbo4calendarius.HappeningSlotTiming{
		Timing: dbo4calendarius.Timing{
			Start: dbo4calendarius.DateTime{Date: "2020-01-01", Time: "10:00"},
			End:   dbo4calendarius.DateTime{Time: "11:00"},
		},
		Repeats: dbo4calendarius.RepeatPeriodOnce,
	}
}

// validHappeningBrief returns a minimal valid HappeningBrief.
func validHappeningBrief() *dbo4calendarius.HappeningBrief {
	return &dbo4calendarius.HappeningBrief{
		HappeningBase: dbo4calendarius.HappeningBase{
			Type:   dbo4calendarius.HappeningTypeSingle,
			Status: dbo4calendarius.HappeningStatusActive,
			Title:  "Test happening",
			Slots: map[string]*dbo4calendarius.HappeningSlot{
				"s1": {HappeningSlotTiming: validSlotTiming()},
			},
		},
	}
}

func TestCreateHappeningRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateHappeningRequest
		wantErr bool
	}{
		{"valid", CreateHappeningRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("space1")},
			Happening:    validHappeningBrief(),
		}, false},
		{"missing_space", CreateHappeningRequest{
			Happening: validHappeningBrief(),
		}, true},
		{"missing_happening", CreateHappeningRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("space1")},
		}, true},
		{"invalid_happening", CreateHappeningRequest{
			SpaceRequest: dto4spaceus.SpaceRequest{SpaceID: coretypes.SpaceID("space1")},
			Happening:    &dbo4calendarius.HappeningBrief{}, // missing type/title/etc.
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateHappeningResponse_Validate(t *testing.T) {
	validDbo := dbo4calendarius.HappeningDbo{
		HappeningBase: dbo4calendarius.HappeningBase{
			Type:   dbo4calendarius.HappeningTypeSingle,
			Status: dbo4calendarius.HappeningStatusActive,
			Title:  "Test happening",
			Slots: map[string]*dbo4calendarius.HappeningSlot{
				"s1": {HappeningSlotTiming: validSlotTiming()},
			},
		},
	}
	validDbo.UserIDs = []string{"u1"}
	validDbo.Dates = []string{"2020-01-01"}
	validDbo.RelatedIDs = []string{"-"}
	tests := []struct {
		name    string
		resp    CreateHappeningResponse
		wantErr bool
	}{
		{"valid", CreateHappeningResponse{ID: "id1", Dbo: validDbo}, false},
		{"missing_id", CreateHappeningResponse{Dbo: validDbo}, true},
		{"invalid_dbo", CreateHappeningResponse{ID: "id1", Dbo: dbo4calendarius.HappeningDbo{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.resp.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
