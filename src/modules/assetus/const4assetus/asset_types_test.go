package const4assetus

import "testing"

func TestValidateAssetType(t *testing.T) {
	tests := []struct {
		name          string
		assetCategory AssetCategory
		assetType     AssetType
		wantErr       bool
	}{
		{"valid_vehicle", AssetCategoryVehicle, AssetTypeVehicleCar, false},
		{"invalid_vehicle_type", AssetCategoryVehicle, "invalid", true},
		{"valid_dwelling", AssetCategoryDwelling, AssetTypeRealEstateHouse, false},
		{"invalid_category", "invalid", AssetTypeVehicleCar, true},
		{"valid_document", AssetCategoryDocument, AssetTypeDocumentTypePassport, false},
		{"valid_sport_gear", AssetCategorySportGear, AssetTypeSportGearBicycle, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateAssetType(tt.assetCategory, tt.assetType); (err != nil) != tt.wantErr {
				t.Errorf("ValidateAssetType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
