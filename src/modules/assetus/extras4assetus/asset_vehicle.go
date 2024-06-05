package extras4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodels/extra"
	"github.com/strongo/validation"
)

func init() {
	RegisterAssetExtraFactory(AssetExtraTypeVehicle, func() extra.Data {
		return new(AssetVehicleExtra)
	})
}

var _ extra.Data = (*AssetVehicleExtra)(nil)

// AssetVehicleExtra is an extension of asset data for vehicles
type AssetVehicleExtra struct {
	extra.BaseData
	WithMakeModelRegNumberFields
	WithEngineData
	Vin string `json:"vin,omitempty" firestore:"vin,omitempty"`
}

func (v AssetVehicleExtra) GetBrief() extra.Data {
	return &AssetVehicleExtra{
		BaseData:                     v.BaseData,
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
	if err := v.BaseData.Validate(); err != nil {
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
