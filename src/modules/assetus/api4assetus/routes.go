package api4assetus

import (
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

// RegisterHttpRoutes registers asset routes
func RegisterHttpRoutes(handle modules.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/assets/create_asset", httpPostCreateAsset)
	handle(http.MethodPost, "/v0/assets/update_asset", httpPostUpdateAsset)
	handle(http.MethodDelete, "/v0/assets/delete_asset", httpDeleteAsset)
}
