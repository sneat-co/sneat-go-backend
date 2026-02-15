package deedus

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/deedus/api4deedus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/deedus/const4deedus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4deedus.ExtensionID, extension.RegisterRoutes(api4deedus.RegisterHttpRoutes))
}
