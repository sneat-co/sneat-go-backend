package const4assetus

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

var AssetRealEstateTypes = []string{
	AssetTypeRealEstateApartment,
	AssetTypeRealEstateHouse,
	AssetTypeRealEstateOffice,
	AssetTypeRealEstateShop,
	AssetTypeRealEstateLand,
	AssetTypeRealEstateGarage,
	AssetTypeRealEstateWarehouse,
}

var AssetVehicleTypes = []string{
	AssetTypeVehicleCar,
	AssetTypeVehicleBus,
	AssetTypeVehicleVan,
	AssetTypeVehicleTruck,
	AssetTypeVehicleMotorcycle,
	AssetTypeVehicleBoat,
	AssetTypeVehicleAircraft,
	AssetTypeVehicleHelicopter,
}

var AssetSportGearTypes = []string{
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

var AssetDocumentTypes = []string{
	AssetTypeDocumentTypePassport,
	AssetTypeDocumentTypeIDCard,
	AssetTypeDocumentTypeDrivingLicense,
	AssetTypeDocumentTypeMarriageCert,
	AssetTypeDocumentTypeBirthCert,
}
