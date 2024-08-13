package api4linkage

import (
	"github.com/sneat-co/sneat-go-core/module"
	"net/http"
)

func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle(http.MethodPost, "/linkage/update_item_relationships", httpUpdateItemRelationships)
}
