package api4listus

import (
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

// RegisterHttpRoutes registers listus routes
func RegisterHttpRoutes(handle modules.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/listus/create_list", httpPostCreateList)
	handle(http.MethodDelete, "/v0/listus/delete_list", httpDeleteList)
	handle(http.MethodPost, "/v0/listus/list_items_create", httpPostCreateListItems)
	handle(http.MethodPost, "/v0/listus/list_items_set_is_done", httpPostSetListItemsIsDone)
	handle(http.MethodDelete, "/v0/listus/list_items_delete", httpDeleteListItems)
	handle(http.MethodPost, "/v0/listus/list_items_reorder", httpPostReorderListItem)
}
