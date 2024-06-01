package models4assetus

import (
	"fmt"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
)

// WithEngineData is a struct for engine data
type WithEngineData struct {
	EngineType EngineType `json:"engineType,omitempty" firestore:"engineType,omitempty"`
	EngineFuel FuelType   `json:"engineFuel,omitempty" firestore:"engineFuel,omitempty"`
	EngineCC   int        `json:"engineCC,omitempty" firestore:"engineCC,omitempty"` // Engine volume in cubic centimetres
	EngineKW   int        `json:"engineKW,omitempty" firestore:"engineKW,omitempty"` // Engine power in kilowatts
	EngineNM   int        `json:"engineNM,omitempty" firestore:"engineNM,omitempty"` // Engine torque in Newton metres
}

// Validate returns error if not valid
func (v WithEngineData) Validate() error {
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

// IsKnownFuelType returns true if a given fuel type is known
func IsKnownFuelType(v FuelType) bool {
	return slice.Contains(FuelTypes, v)
}

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
