package models4assetus

import (
	"fmt"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
	"strings"
)

// EngineType is a type of engine
type EngineType = string

const (
	EngineTypeUnknown    EngineType = ""
	EngineTypeOther      EngineType = "other"
	EngineTypeCombustion EngineType = "combustion"
	EngineTypeElectric   EngineType = "electric"
	EngineTypePHEV       EngineType = "phev"
	EngineTypeHybrid     EngineType = "hybrid"
	EngineTypeSteam      EngineType = "steam"
)

// FuelType is a type of a fuel
type FuelType string

const (
	FuelTypeUnknown  FuelType = ""
	FuelTypeOther    FuelType = "other"
	FuelTypeBio      FuelType = "bio"
	FuelTypePetrol   FuelType = "petrol"
	FuelTypeDiesel   FuelType = "diesel"
	FuelTypeHydrogen FuelType = "hydrogen"
)

// FuelTypes is a list of known fuel types
var FuelTypes = []FuelType{
	FuelTypeUnknown,
	FuelTypeOther,
	FuelTypePetrol,
	FuelTypeDiesel,
	FuelTypeHydrogen,
	FuelTypeBio,
}

// IsKnownFuelType returns true if given fuel type is known
func IsKnownFuelType(v FuelType) bool {
	return slice.Contains(FuelTypes, v)
}

// VehicleData is extension of asset data for vehicles
type VehicleData struct {
	EngineData
	Vin string `json:"vin,omitempty" firestore:"vin"`
}

// EngineData is a struct for engine data
type EngineData struct {
	EngineType EngineType `json:"engineType" firestore:"engineType"`
	EngineFuel FuelType   `json:"engineFuel" firestore:"engineFuel"`
	EngineCC   int        `json:"engineCC" firestore:"engineCC"` // Engine volume in cubic centimetres
	EngineKW   int        `json:"engineKW" firestore:"engineKW"` // Engine power in kilowatts
	EngineNM   int        `json:"engineNM" firestore:"engineNM"` // Engine torque in Newton metres
}

// Validate returns error if not valid
func (v EngineData) Validate() error {
	switch v.EngineType {
	case EngineTypeUnknown, EngineTypeOther:
		if !IsKnownFuelType(v.EngineFuel) {
			return validation.NewErrBadRecordFieldValue("fuelType", fmt.Sprintf("unknown fuel type: %s", v.EngineFuel))
		}
	case EngineTypeCombustion, EngineTypeHybrid, EngineTypePHEV:
		switch v.EngineFuel {
		case FuelTypePetrol, FuelTypeDiesel, FuelTypeHydrogen, FuelTypeUnknown, FuelTypeOther:
		//OK
		default:
			return validation.NewErrBadRecordFieldValue("fuelType", fmt.Sprintf("unknown fuel type: %s", v.EngineFuel))
		}
	case EngineTypeElectric:
		switch v.EngineFuel {
		case FuelTypeUnknown, FuelTypeOther:
		//OK
		default:
			return validation.NewErrBadRecordFieldValue("fuelType", fmt.Sprintf("unknown fuel type: %s", v.EngineFuel))
		}
	case EngineTypeSteam:
		switch v.EngineFuel {
		case FuelTypeUnknown, FuelTypeOther:
		//OK
		default:
			return validation.NewErrBadRecordFieldValue("fuelType", fmt.Sprintf("unknown fuel type: %s", v.EngineFuel))
		}
	default:
		return validation.NewErrBadRecordFieldValue("engineType", "unknown engine type: "+v.EngineType)
	}
	return nil
}

var _ AssetMain = (*VehicleAssetMainData)(nil)

// VehicleAssetMainData is a base struct for DB and request DTOs
type VehicleAssetMainData struct {
	AssetMainDto
	VehicleData
}

// SpecificData returns specific data
func (v *VehicleAssetMainData) SpecificData() AssetSpecificData {
	return &v.VehicleData
}

// SetSpecificData sets specific data
func (v *VehicleAssetMainData) SetSpecificData(data AssetSpecificData) {
	v.VehicleData = data.(VehicleData)
}

// GenerateTitle generates asset title from vehicle data
func (v *VehicleAssetMainData) GenerateTitle() string {
	if v.Title != "" {
		return v.Title
	}
	title := make([]string, 0, 2)
	if v.Make != "" {
		title = append(title, v.Make)
	}
	if v.Model != "" {
		title = append(title, v.Model)
	}
	if v.RegNumber != "" {
		title = append(title, v.RegNumber)
	}
	return strings.Join(title, " ")
}

// Validate returns error if not valid
func (v *VehicleAssetMainData) Validate() error {
	return nil
}

var _ AssetDbData = (*AssetDtoVehicle)(nil)

// AssetDtoVehicle is a DB DTO
type AssetDtoVehicle struct {
	VehicleAssetMainData
	AssetExtraDto
}

// NewVehicleAssetDbData creates new AssetDtoVehicle
func NewVehicleAssetDbData() *AssetDtoVehicle {
	return &AssetDtoVehicle{
		//VehicleAssetMainData: new(VehicleAssetMainData),
		//AssetExtraDto:        new(AssetExtraDto),
	}
}

// Validate returns error if not valid
func (v *AssetDtoVehicle) Validate() error {
	if err := v.VehicleAssetMainData.Validate(); err != nil {
		return err
	}
	if err := v.AssetExtraDto.Validate(); err != nil {
		return err
	}
	return nil
}
