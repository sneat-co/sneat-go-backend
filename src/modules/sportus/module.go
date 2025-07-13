package sportus

import (
	entities2 "github.com/sneat-co/sneat-core-modules/generic/entities"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/api4sportus"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/const4sportus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	entities2.Register(
		entities2.Entity{Name: "Spot", AllowCreate: true, AllowUpdate: true},
	)
	return module.NewExtension(const4sportus.ModuleID, module.RegisterRoutes(api4sportus.RegisterRoutes))
}
