package sportus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/generic/entities"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/api4sportus"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus/const4sportus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	entities.Register(
		entities.Entity{Name: "Spot", AllowCreate: true, AllowUpdate: true},
	)
	return modules.NewModule(const4sportus.ModuleID, api4sportus.RegisterRoutes)
}
