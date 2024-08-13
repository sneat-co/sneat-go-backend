package api4sportus

import (
	"github.com/sneat-co/sneat-go-core/module"
)

// RegisterRoutes registers HTTP handle
func RegisterRoutes(handle module.HTTPHandleFunc) {
	if handle == nil {
		panic("handle == nil")
	}
	registerSpotHandlers(handle)
	registerQuiverHandlers(handle)
}
