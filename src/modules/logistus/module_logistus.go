package logistus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/api4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/const4logistus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4logistus.ModuleID, api4logist.RegisterHttpRoutes)
}
