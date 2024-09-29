package extras4assetus

import (
	extra2 "github.com/sneat-co/sneat-core-modules/core/extra"
)

const (
	AssetExtraTypeVehicle  extra2.Type = "vehicle"
	AssetExtraTypeDwelling extra2.Type = "dwelling"
	AssetExtraTypeDocument extra2.Type = "document"
)

func init() {
	extra2.RegisterFactory(AssetExtraTypeVehicle, func() extra2.Data { return &AssetVehicleExtra{} })
	extra2.RegisterFactory(AssetExtraTypeDwelling, func() extra2.Data { return &AssetDwellingExtra{} })
	extra2.RegisterFactory(AssetExtraTypeDocument, func() extra2.Data { return &AssetDocumentExtra{} })
}
