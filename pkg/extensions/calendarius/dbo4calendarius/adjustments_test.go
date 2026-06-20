package dbo4calendarius

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func TestHappeningAdjustment_IsEmpty(t *testing.T) {
	var nilAdj *HappeningAdjustment
	assert.True(t, nilAdj.IsEmpty())
	assert.True(t, (&HappeningAdjustment{}).IsEmpty())
	assert.False(t, (&HappeningAdjustment{Slots: map[string]*SlotAdjustment{"s1": {}}}).IsEmpty())
}

func TestHappeningAdjustment_Validate(t *testing.T) {
	var nilAdj *HappeningAdjustment
	assert.NoError(t, nilAdj.Validate())

	assert.NoError(t, (&HappeningAdjustment{Slots: map[string]*SlotAdjustment{"s1": {}}}).Validate())

	badCancellation := &Cancellation{Reason: string(make([]byte, ReasonMaxLen+1))}
	bad := &HappeningAdjustment{Slots: map[string]*SlotAdjustment{
		"s1": {Cancellation: badCancellation},
	}}
	assert.Error(t, bad.Validate())
}

func TestSlotAdjustment_IsEmpty(t *testing.T) {
	assert.True(t, (&SlotAdjustment{}).IsEmpty())

	withCancellation := &SlotAdjustment{Cancellation: &Cancellation{}}
	assert.False(t, withCancellation.IsEmpty())

	withAdjustment := &SlotAdjustment{
		Adjustment: &HappeningSlot{Locations: []Location{{Type: "online", Title: "Zoom"}}},
	}
	assert.False(t, withAdjustment.IsEmpty())
}

func TestSlotAdjustment_Validate(t *testing.T) {
	var nilAdj *SlotAdjustment
	assert.Error(t, nilAdj.Validate())

	assert.NoError(t, (&SlotAdjustment{}).Validate())

	validCanc := &Cancellation{At: time.Now(), By: dbmodels.ByUser{UID: "u1"}}
	assert.NoError(t, (&SlotAdjustment{Cancellation: validCanc}).Validate())

	badCanc := &Cancellation{}
	assert.Error(t, (&SlotAdjustment{Cancellation: badCanc}).Validate())
}

func TestCancellation_Validate(t *testing.T) {
	tests := []struct {
		name    string
		c       Cancellation
		wantErr bool
	}{
		{"valid", Cancellation{At: time.Now(), By: dbmodels.ByUser{UID: "u1"}}, false},
		{"missing_at", Cancellation{By: dbmodels.ByUser{UID: "u1"}}, true},
		{"missing_by", Cancellation{At: time.Now()}, true},
		{"reason_too_long", Cancellation{At: time.Now(), By: dbmodels.ByUser{UID: "u1"}, Reason: string(make([]byte, ReasonMaxLen+1))}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
