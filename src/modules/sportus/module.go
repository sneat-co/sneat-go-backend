package sportus

import (
	entities2 "github.com/sneat-co/sneat-core-modules/generic/entities"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/api4sportus"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/const4sportus"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	entities2.Register(
		entities2.Entity{Name: "Spot", AllowCreate: true, AllowUpdate: true},
	)
	return extension.NewExtension(const4sportus.ExtensionID,
		extension.RegisterRoutes(api4sportus.RegisterRoutes),
	)
}
