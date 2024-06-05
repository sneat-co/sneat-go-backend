package extras4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodels/extra"
)

const (
	AssetExtraTypeVehicle  extra.Type = "vehicle"
	AssetExtraTypeDwelling extra.Type = "dwelling"
	AssetExtraTypeDocument extra.Type = "document"
)

func init() {
	extra.RegisterFactory(AssetExtraTypeVehicle, func() extra.Data { return &AssetVehicleExtra{} })
	extra.RegisterFactory(AssetExtraTypeDwelling, func() extra.Data { return &AssetDwellingExtra{} })
	extra.RegisterFactory(AssetExtraTypeDocument, func() extra.Data { return &AssetDocumentExtra{} })
}
