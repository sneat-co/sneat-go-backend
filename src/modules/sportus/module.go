package sportus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/generic/entities"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/api4sportus"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/const4sportus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	entities.Register(
		entities.Entity{Name: "Spot", AllowCreate: true, AllowUpdate: true},
	)
	return module.NewModule(const4sportus.ModuleID, module.RegisterRoutes(api4sportus.RegisterRoutes))
}
