package logistus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/api4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/const4logistus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4logistus.ModuleID, module.RegisterRoutes(api4logist.RegisterHttpRoutes))
}
