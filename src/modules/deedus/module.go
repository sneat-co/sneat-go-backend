package deedus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/deedus/api4deedus"
	"github.com/sneat-co/sneat-go-backend/src/modules/deedus/const4deedus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4deedus.ExtensionID, extension.RegisterRoutes(api4deedus.RegisterHttpRoutes))
}
