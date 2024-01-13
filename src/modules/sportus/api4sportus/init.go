package api4sportus

import (
	"github.com/sneat-co/sneat-go-core/modules"
)

// RegisterRoutes registers HTTP handle
func RegisterRoutes(handle modules.HTTPHandleFunc) {
	if handle == nil {
		panic("handle == nil")
	}
	registerSpotHandlers(handle)
	registerQuiverHandlers(handle)
}
