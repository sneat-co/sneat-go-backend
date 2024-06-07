package const4assetus

import (
	"fmt"
	"github.com/strongo/slice"
	"github.com/strongo/validation"
)

// AssetType is a type of asset, e.g. house, apartment, car, boat, etc
type AssetType = string

const (
	AssetTypeSportGearBicycle = "bicycle"

	AssetTypeSportGearKite          = "kite"
	AssetTypeSportGearKiteBar       = "kite_bar"
	AssetTypeSportGearKiteBoard     = "kite_board"
	AssetTypeSportGearKiteHydrofoil = "kite_hydrofoil"

	AssetTypeSportGearProneHydrofoil = "prone_hydrofoil"

	AssetTypeSportGearSurfBoard = "surf_board"

	AssetTypeSportGearWetsuit = "wetsuit"

	AssetTypeSportGearWing          = "wing"
	AssetTypeSportGearWingBoard     = "wing_board"
	AssetTypeSportGearWingHydrofoil = "wing_hydrofoil"
)

const (
	AssetTypeVehicleAircraft   = "aircraft"
	AssetTypeVehicleBoat       = "boat"
	AssetTypeVehicleBus        = "bus"
	AssetTypeVehicleCar        = "car"
	AssetTypeVehicleHelicopter = "helicopter"
	AssetTypeVehicleMotorcycle = "motorcycle"
	AssetTypeVehicleTruck      = "truck"
	AssetTypeVehicleVan        = "van"
)

const (
	AssetTypeRealEstateApartment = "apartment"
	AssetTypeRealEstateHouse     = "house"
	AssetTypeRealEstateOffice    = "office"
	AssetTypeRealEstateShop      = "shop"
	AssetTypeRealEstateLand      = "land"
	AssetTypeRealEstateGarage    = "garage"
	AssetTypeRealEstateWarehouse = "warehouse"
)

var DwellingAssetTypes = []string{
	AssetTypeRealEstateApartment,
	AssetTypeRealEstateHouse,
	AssetTypeRealEstateOffice,
	AssetTypeRealEstateShop,
	AssetTypeRealEstateLand,
	AssetTypeRealEstateGarage,
	AssetTypeRealEstateWarehouse,
}

var VehicleAssetTypes = []string{
	AssetTypeVehicleCar,
	AssetTypeVehicleBus,
	AssetTypeVehicleVan,
	AssetTypeVehicleTruck,
	AssetTypeVehicleMotorcycle,
	AssetTypeVehicleBoat,
	AssetTypeVehicleAircraft,
	AssetTypeVehicleHelicopter,
}

var SportGearAssetTypes = []string{
	AssetTypeSportGearBicycle,
	AssetTypeSportGearKite,
	AssetTypeSportGearKiteBar,
	AssetTypeSportGearKiteBoard,
	AssetTypeSportGearKiteHydrofoil,
	AssetTypeSportGearProneHydrofoil,
	AssetTypeSportGearSurfBoard,
	AssetTypeSportGearWetsuit,
	AssetTypeSportGearWing,
	AssetTypeSportGearWingBoard,
	AssetTypeSportGearWingHydrofoil,
}

const (
	AssetTypeDocumentTypePassport       = "passport"
	AssetTypeDocumentTypeIDCard         = "id_card"
	AssetTypeDocumentTypeDrivingLicense = "driving_license"
	AssetTypeDocumentTypeMarriageCert   = "marriage_cert"
	AssetTypeDocumentTypeBirthCert      = "birth_cert"
)

var DocumentAssetTypes = []string{
	AssetTypeDocumentTypePassport,
	AssetTypeDocumentTypeIDCard,
	AssetTypeDocumentTypeDrivingLicense,
	AssetTypeDocumentTypeMarriageCert,
	AssetTypeDocumentTypeBirthCert,
}

var assetTypesByCategory = map[AssetCategory][]string{
	AssetCategoryVehicle:   VehicleAssetTypes,
	AssetCategoryDwelling:  DwellingAssetTypes,
	AssetCategorySportGear: SportGearAssetTypes,
	AssetCategoryDocument:  DocumentAssetTypes,
}

func ValidateAssetType(assetCategory AssetCategory, assetType AssetType) error {
	if types, ok := assetTypesByCategory[assetCategory]; ok {
		if !slice.Contains(types, assetType) {
			return validation.NewErrBadRecordFieldValue("type", fmt.Sprintf("unknown %s type: %s", assetCategory, assetType))
		}
		return nil
	}
	return validation.NewErrBadRecordFieldValue("assetCategory", "unknown value: "+string(assetCategory))
}
