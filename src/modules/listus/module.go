package listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/api4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4listus.ExtensionID,
		extension.RegisterRoutes(api4listus.RegisterHttpRoutes),
	)
}
