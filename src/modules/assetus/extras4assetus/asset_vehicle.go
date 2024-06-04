package extras4assetus

import (
	"github.com/strongo/validation"
)

func init() {
	RegisterAssetExtraFactory(AssetExtraTypeVehicle, func() AssetExtra {
		return new(AssetVehicleExtra)
	})
}

var _ AssetExtra = (*AssetVehicleExtra)(nil)

// AssetVehicleExtra is an extension of asset data for vehicles
type AssetVehicleExtra struct {
	AssetExtraBase
	WithMakeModelRegNumberFields
	WithEngineData
	Vin string `json:"vin,omitempty" firestore:"vin,omitempty"`
}

func (v AssetVehicleExtra) GetBrief() AssetExtra {
	return &AssetVehicleExtra{
		AssetExtraBase:               v.AssetExtraBase,
		WithMakeModelRegNumberFields: v.WithMakeModelRegNumberFields,
		Vin:                          v.Vin,
	}
}

func (v AssetVehicleExtra) RequiredFields() []string {
	return nil
}

func (v AssetVehicleExtra) IndexedFields() []string {
	return []string{"make", "model", "make+model", "regNumber", "vin"}
}

func (v AssetVehicleExtra) Validate() error {
	if err := v.AssetExtraBase.Validate(); err != nil {
		return err
	}
	if err := v.WithMakeModelRegNumberFields.Validate(); err != nil {
		return err
	}
	if err := v.WithEngineData.Validate(); err != nil {
		return validation.NewErrBadRecordFieldValue("engineData", err.Error())
	}
	return nil
}
