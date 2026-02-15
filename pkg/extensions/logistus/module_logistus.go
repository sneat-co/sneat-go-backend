package logistus

import (
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/api4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/const4logistus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4logistus.ExtensionID,
		extension.RegisterRoutes(api4logist.RegisterHttpRoutes),
	)
}
