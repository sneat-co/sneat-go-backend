package api4assetus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-core/extension"
)

// RegisterHttpRoutes registers asset routes
func RegisterHttpRoutes(handle extension.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/assets/create_asset", httpPostCreateAsset)
	handle(http.MethodPost, "/v0/assets/update_asset", httpPostUpdateAsset)
	handle(http.MethodDelete, "/v0/assets/delete_asset", httpDeleteAsset)
	handle(http.MethodPost, "/v0/assets/create_vehicle_record", httpPostCreateVehicleRecord)
}
