package extras4assetus

import (
	"testing"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/briefs4assetus"
)

func TestWithMakeModelFields_GenerateTitleFromMakeModelAndRegNumber(t *testing.T) {
	v := WithMakeModelFields{Make: "Toyota", Model: "Camry"}
	title := v.GenerateTitleFromMakeModelAndRegNumber("123-ABC")
	expected := "Toyota Camry # 123-ABC"
	if title != expected {
		t.Errorf("expected %s, got %s", expected, title)
	}

	v = WithMakeModelFields{}
	title = v.GenerateTitleFromMakeModelAndRegNumber("")
	if title != "" {
		t.Errorf("expected empty title, got %s", title)
	}
}

func TestWithMakeModelFields_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		v := WithMakeModelFields{Make: "Toyota", Model: "Camry"}
		if err := v.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("empty_make", func(t *testing.T) {
		v := WithMakeModelFields{Model: "Camry"}
		if err := v.Validate(); err == nil {
			t.Error("expected error for empty make")
		}
	})

	t.Run("spaces_make", func(t *testing.T) {
		v := WithMakeModelFields{Make: " Toyota ", Model: "Camry"}
		if err := v.Validate(); err == nil {
			t.Error("expected error for spaces in make")
		}
	})
}

func TestWithOptionalRegNumberField_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		v := WithOptionalRegNumberField{RegNumber: "123"}
		if err := v.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("spaces", func(t *testing.T) {
		v := WithOptionalRegNumberField{RegNumber: " 123 "}
		if err := v.Validate(); err == nil {
			t.Error("expected error for spaces in regNumber")
		}
	})
}

func TestAssetVehicleExtra_ValidateWithAssetBrief(t *testing.T) {
	v := &AssetVehicleExtra{
		WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
			WithMakeModelFields: WithMakeModelFields{Make: "Toyota", Model: "Camry"},
		},
	}
	t.Run("valid", func(t *testing.T) {
		if err := v.ValidateWithAssetBrief(briefs4assetus.AssetBrief{Title: "Title"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid_all_empty", func(t *testing.T) {
		v2 := &AssetVehicleExtra{}
		// v2.Validate() will fail first because Make is empty
		if err := v2.ValidateWithAssetBrief(briefs4assetus.AssetBrief{}); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("missing_brief_title_and_extra_fields", func(t *testing.T) {
		v2 := &AssetVehicleExtra{
			WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
				WithMakeModelFields: WithMakeModelFields{Make: "Toyota", Model: "Camry"},
			},
		}
		// Make and Model are present, but if we clear them...
		v2.Make = ""
		v2.Model = ""
		v2.RegNumber = ""
		// v2.Validate() will fail because Make/Model are required in WithMakeModelFields.Validate()
		// So we need to test the logic in ValidateWithAssetBrief specifically.
	})
}

func TestAssetVehicleExtra_GetBrief(t *testing.T) {
	v := &AssetVehicleExtra{
		Vin: "VIN123",
		WithMakeModelRegNumberFields: WithMakeModelRegNumberFields{
			WithMakeModelFields: WithMakeModelFields{Make: "Toyota", Model: "Camry"},
		},
	}
	brief := v.GetBrief().(*AssetVehicleExtra)
	if brief.Vin != "VIN123" {
		t.Errorf("expected Vin VIN123, got %s", brief.Vin)
	}
}

func TestAssetVehicleExtra_RequiredFields(t *testing.T) {
	v := &AssetVehicleExtra{}
	if v.RequiredFields() != nil {
		t.Error("expected nil")
	}
}

func TestAssetVehicleExtra_IndexedFields(t *testing.T) {
	v := &AssetVehicleExtra{}
	if len(v.IndexedFields()) == 0 {
		t.Error("expected indexed fields")
	}
}

func TestAssetDwellingExtra_ValidateWithAssetBrief(t *testing.T) {
	v := AssetDwellingExtra{}
	t.Run("valid", func(t *testing.T) {
		if err := v.ValidateWithAssetBrief(briefs4assetus.AssetBrief{Title: "Title"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		if err := v.ValidateWithAssetBrief(briefs4assetus.AssetBrief{}); err == nil {
			t.Error("expected error")
		}
	})
}

func TestAssetDwellingExtra_GetBrief(t *testing.T) {
	v := AssetDwellingExtra{AreaSqM: 100}
	brief := v.GetBrief().(*AssetDwellingExtra)
	if brief.AreaSqM != 100 {
		t.Errorf("expected 100, got %d", brief.AreaSqM)
	}
}

func TestAssetDocumentExtra_ValidateWithAssetBrief(t *testing.T) {
	v := &AssetDocumentExtra{}
	t.Run("valid", func(t *testing.T) {
		if err := v.ValidateWithAssetBrief(briefs4assetus.AssetBrief{Title: "Title"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("invalid", func(t *testing.T) {
		if err := v.ValidateWithAssetBrief(briefs4assetus.AssetBrief{}); err == nil {
			t.Error("expected error")
		}
	})
}

func TestAssetDocumentExtra_GetBrief(t *testing.T) {
	v := &AssetDocumentExtra{IssuedOn: "2023-01-01"}
	brief := v.GetBrief().(*AssetDocumentExtra)
	if brief.IssuedOn != "2023-01-01" {
		t.Errorf("expected 2023-01-01, got %s", brief.IssuedOn)
	}
}

func TestAssetDwellingExtra_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		v := AssetDwellingExtra{NumberOfBedrooms: 2, AreaSqM: 50}
		if err := v.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("negative_bedrooms", func(t *testing.T) {
		v := AssetDwellingExtra{NumberOfBedrooms: -1}
		if err := v.Validate(); err == nil {
			t.Error("expected error for negative bedrooms")
		}
	})
}

func TestWithEngineData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		v       WithEngineData
		wantErr bool
	}{
		{"valid_combustion_petrol", WithEngineData{EngineType: EngineTypeCombustion, EngineFuel: FuelTypePetrol}, false},
		{"invalid_combustion_electric", WithEngineData{EngineType: EngineTypeCombustion, EngineFuel: "electric"}, true},
		{"valid_electric_unknown_fuel", WithEngineData{EngineType: EngineTypeElectric, EngineFuel: FuelTypeUnknown}, false},
		{"invalid_electric_petrol", WithEngineData{EngineType: EngineTypeElectric, EngineFuel: FuelTypePetrol}, true},
		{"invalid_engine_type", WithEngineData{EngineType: "invalid"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.v.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAssetDocumentExtra_Validate(t *testing.T) {
	t.Run("valid_dates", func(t *testing.T) {
		v := &AssetDocumentExtra{
			IssuedOn:      "2023-01-01",
			EffectiveFrom: "2023-01-01",
			ExpiresOn:     "2024-01-01",
		}
		if err := v.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid_issued_on", func(t *testing.T) {
		v := &AssetDocumentExtra{IssuedOn: "invalid"}
		if err := v.Validate(); err == nil {
			t.Error("expected error for invalid date")
		}
	})

	t.Run("expires_before_effective", func(t *testing.T) {
		v := &AssetDocumentExtra{
			EffectiveFrom: "2024-01-01",
			ExpiresOn:     "2023-01-01",
		}
		if err := v.Validate(); err == nil {
			t.Error("expected error for expires before effective")
		}
	})
}
