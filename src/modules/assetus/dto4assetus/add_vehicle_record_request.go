package dto4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/validate"
	"github.com/strongo/validation"
)

type AddVehicleRecordRequest struct {
	AssetID        string
	FuelVolume     float32
	FuelVolumeUnit string
	FuelCost       float32
	Currency       string
	Mileage        float32
	MileageUnit    string
	dto4teamus.TeamRequest
}

func (v AddVehicleRecordRequest) Validate() error {
	if err := v.TeamRequest.Validate(); err != nil {
		return err
	}
	if err := validate.RecordID(v.AssetID); err != nil {
		return validation.NewErrBadRequestFieldValue("assetID", err.Error())
	}
	return nil
}
