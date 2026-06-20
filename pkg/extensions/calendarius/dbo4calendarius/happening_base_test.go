package dbo4calendarius

import (
	"testing"
	"time"

	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
)

func validCancellation() Cancellation {
	return Cancellation{
		At: time.Now(),
		By: dbmodels.ByUser{UID: "u1"},
	}
}

func validHappeningBase() *HappeningBase {
	return &HappeningBase{
		Type:   HappeningTypeSingle,
		Status: HappeningStatusActive,
		Title:  "Test happening",
		Slots: map[string]*HappeningSlot{
			"s1": {HappeningSlotTiming: validSlotTiming()},
		},
	}
}

func TestHappeningBase_MarkAsCanceled(t *testing.T) {
	v := validHappeningBase()
	updates := v.MarkAsCanceled(validCancellation())
	assert.Equal(t, HappeningStatusCanceled, v.Status)
	assert.NotNil(t, v.Cancellation)
	assert.Len(t, updates, 2)

	// Idempotent: already canceled produces no new updates.
	v2 := v.MarkAsCanceled(validCancellation())
	assert.Empty(t, v2)
}

func TestHappeningBase_RemoveCancellation(t *testing.T) {
	v := validHappeningBase()
	v.MarkAsCanceled(validCancellation())

	updates := v.RemoveCancellation()
	assert.Equal(t, HappeningStatusActive, v.Status)
	assert.Nil(t, v.Cancellation)
	assert.Len(t, updates, 2)

	// Idempotent: nothing to remove.
	updates2 := v.RemoveCancellation()
	assert.Empty(t, updates2)
}

func TestHappeningBase_GetSlot_HasSlot(t *testing.T) {
	v := validHappeningBase()
	assert.True(t, v.HasSlot("s1"))
	assert.NotNil(t, v.GetSlot("s1"))
	assert.False(t, v.HasSlot("missing"))
	assert.Nil(t, v.GetSlot("missing"))
}

func TestHappeningBase_Validate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(v *HappeningBase)
		wantErr bool
	}{
		{"valid", func(v *HappeningBase) {}, false},
		{"empty_type", func(v *HappeningBase) { v.Type = "" }, true},
		{"unknown_type", func(v *HappeningBase) { v.Type = "weird" }, true},
		{"empty_title", func(v *HappeningBase) { v.Title = "" }, true},
		{"title_too_long", func(v *HappeningBase) { v.Title = string(make([]byte, 101)) }, true},
		{"summary_too_long", func(v *HappeningBase) { v.Summary = string(make([]byte, 201)) }, true},
		{"empty_status", func(v *HappeningBase) { v.Status = "" }, true},
		{"unknown_status", func(v *HappeningBase) { v.Status = "weird" }, true},
		{"canceled_without_cancellation", func(v *HappeningBase) { v.Status = HappeningStatusCanceled }, true},
		{
			name: "cancellation_on_active",
			modify: func(v *HappeningBase) {
				c := validCancellation()
				v.Cancellation = &c
			},
			wantErr: true,
		},
		{"no_slots", func(v *HappeningBase) { v.Slots = nil }, true},
		{
			name: "empty_slot_key",
			modify: func(v *HappeningBase) {
				v.Slots = map[string]*HappeningSlot{"  ": {HappeningSlotTiming: validSlotTiming()}}
			},
			wantErr: true,
		},
		{
			name: "invalid_slot",
			modify: func(v *HappeningBase) {
				v.Slots = map[string]*HappeningSlot{"s1": {HappeningSlotTiming: HappeningSlotTiming{Repeats: ""}}}
			},
			wantErr: true,
		},
		{
			name: "price_with_empty_id",
			modify: func(v *HappeningBase) {
				v.Prices = []*HappeningPrice{{Term: Term{Unit: TermUnitHour, Length: 1}, Amount: validAmount()}}
			},
			wantErr: true,
		},
		{
			name: "valid_price",
			modify: func(v *HappeningBase) {
				v.Prices = []*HappeningPrice{{ID: "p1", Term: Term{Unit: TermUnitHour, Length: 1}, Amount: validAmount()}}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validHappeningBase()
			tt.modify(v)
			err := v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
