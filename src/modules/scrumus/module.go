package scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/api4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/const4srumus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	return extension.NewExtension(const4srumus.ExtensionID,
		extension.RegisterRoutes(api4scrumus.RegisterHttpRoutes),
	)
}
