package api4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/extras4assetus"
	"testing"
)

func TestCreateAssetBaseDbo(t *testing.T) {
	tests := []struct {
		category string
		wantErr  bool
	}{
		{string(extras4assetus.AssetExtraTypeVehicle), false},
		{string(extras4assetus.AssetExtraTypeDwelling), false},
		{string(extras4assetus.AssetExtraTypeDocument), false},
		{"invalid", true},
		{"", true},
	}
	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			_, err := createAssetBaseDbo(tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("createAssetBaseDbo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
