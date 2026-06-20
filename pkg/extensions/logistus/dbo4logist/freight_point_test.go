package dbo4logist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFreightPoint_HasTask(t *testing.T) {
	v := FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskPick, ShippingPointTaskLoad}}
	assert.True(t, v.HasTask(ShippingPointTaskPick))
	assert.True(t, v.HasTask(ShippingPointTaskLoad))
	assert.False(t, v.HasTask(ShippingPointTaskUnload))
}

func TestFreightPoint_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       FreightPoint
		wantErr bool
	}{
		{"valid_pick", FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskPick}}, false},
		{"empty_tasks", FreightPoint{}, true},
		{"valid_load", FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskLoad}, ToLoad: &FreightLoad{NumberOfPallets: 1}}, false},
		{"load_without_task", FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskPick}, ToLoad: &FreightLoad{NumberOfPallets: 1}}, true},
		{"unload_without_task", FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskPick}, ToUnload: &FreightLoad{NumberOfPallets: 1}}, true},
		{"bad_load", FreightPoint{Tasks: []ShippingPointTask{ShippingPointTaskLoad}, ToLoad: &FreightLoad{NumberOfPallets: -1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
