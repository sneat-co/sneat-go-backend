package api4generic

import (
	"github.com/sneat-co/sneat-go-core/module"
)

// RegisterHttpRoutes registers HTTP handlers
func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle("POST", "/api4invitus/$generic/create", create)
	handle("PUT", "/api4invitus/$generic/update", update)
	handle("DELETE", "/api4invitus/$generic/delete", delete)
}
