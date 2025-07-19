package logistus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/api4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/const4logistus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4logistus.ExtensionID,
		extension.RegisterRoutes(api4logist.RegisterHttpRoutes),
	)
}
