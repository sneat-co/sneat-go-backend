package api4generic

import (
	"github.com/sneat-co/sneat-go-core/module"
	"net/http"
)

// RegisterHttpRoutes registers HTTP handlers
func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle(http.MethodPost, "/api4invitus/$generic/create", create)
	handle(http.MethodPut, "/api4invitus/$generic/update", update)
	handle(http.MethodDelete, "/api4invitus/$generic/delete", delete)
}
