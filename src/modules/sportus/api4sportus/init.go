package api4sportus

import (
	"github.com/sneat-co/sneat-go-core/extension"
)

// RegisterRoutes registers HTTP handle
func RegisterRoutes(handle extension.HTTPHandleFunc) {
	if handle == nil {
		panic("handle == nil")
	}
	registerSpotHandlers(handle)
	registerQuiverHandlers(handle)
}
