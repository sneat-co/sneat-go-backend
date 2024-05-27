package api4linkage

import (
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

func RegisterHttpRoutes(handle modules.HTTPHandleFunc) {
	handle(http.MethodPost, "/linkage/update_item_relationships", httpUpdateItemRelationships)
}
