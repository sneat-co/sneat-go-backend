package dto4calendarius

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/stretchr/testify/assert"
)

func validSlotWithID() HappeningSlotWithID {
	return HappeningSlotWithID{
		ID:            "slot1",
		HappeningSlot: dbo4calendarius.HappeningSlot{HappeningSlotTiming: validSlotTiming()},
	}
}

func TestHappeningSlotWithID_Validate(t *testing.T) {
	tests := []struct {
		name    string
		slot    HappeningSlotWithID
		wantErr bool
	}{
		{"valid", validSlotWithID(), false},
		{"missing_id", HappeningSlotWithID{
			HappeningSlot: dbo4calendarius.HappeningSlot{HappeningSlotTiming: validSlotTiming()},
		}, true},
		{"invalid_slot", HappeningSlotWithID{
			ID:            "slot1",
			HappeningSlot: dbo4calendarius.HappeningSlot{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.slot.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHappeningSlotRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     HappeningSlotRequest
		wantErr bool
	}{
		{"valid", HappeningSlotRequest{
			HappeningRequest: validHappeningRequest(),
			Slot:             validSlotWithID(),
		}, false},
		{"invalid_happening_request", HappeningSlotRequest{
			HappeningRequest: HappeningRequest{},
			Slot:             validSlotWithID(),
		}, true},
		{"invalid_slot", HappeningSlotRequest{
			HappeningRequest: validHappeningRequest(),
			Slot:             HappeningSlotWithID{},
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

func TestHappeningSlotRefRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     HappeningSlotRefRequest
		wantErr bool
	}{
		{"valid", HappeningSlotRefRequest{
			HappeningRequest: validHappeningRequest(),
			SlotID:           "slot1",
		}, false},
		{"invalid_happening_request", HappeningSlotRefRequest{
			HappeningRequest: HappeningRequest{},
			SlotID:           "slot1",
		}, true},
		{"missing_slot_id", HappeningSlotRefRequest{
			HappeningRequest: validHappeningRequest(),
		}, true},
		{"blank_slot_id", HappeningSlotRefRequest{
			HappeningRequest: validHappeningRequest(),
			SlotID:           "   ",
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

func TestDeleteHappeningSlotRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     DeleteHappeningSlotRequest
		wantErr bool
	}{
		{"valid", DeleteHappeningSlotRequest{
			HappeningSlotRefRequest: HappeningSlotRefRequest{
				HappeningRequest: validHappeningRequest(),
				SlotID:           "slot1",
			},
		}, false},
		{"missing_slot_id", DeleteHappeningSlotRequest{
			HappeningSlotRefRequest: HappeningSlotRefRequest{
				HappeningRequest: validHappeningRequest(),
			},
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

func TestHappeningSlotDateRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     HappeningSlotDateRequest
		wantErr bool
	}{
		{"valid", HappeningSlotDateRequest{
			HappeningRequest: validHappeningRequest(),
			Date:             "2020-01-01",
			Slot:             validSlotWithID(),
		}, false},
		{"invalid_happening_request", HappeningSlotDateRequest{
			HappeningRequest: HappeningRequest{},
			Date:             "2020-01-01",
			Slot:             validSlotWithID(),
		}, true},
		{"missing_date", HappeningSlotDateRequest{
			HappeningRequest: validHappeningRequest(),
			Slot:             validSlotWithID(),
		}, true},
		{"bad_date", HappeningSlotDateRequest{
			HappeningRequest: validHappeningRequest(),
			Date:             "not-a-date",
			Slot:             validSlotWithID(),
		}, true},
		{"invalid_slot", HappeningSlotDateRequest{
			HappeningRequest: validHappeningRequest(),
			Date:             "2020-01-01",
			Slot:             HappeningSlotWithID{},
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
