package debtus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4debtus.ModuleID,
		module.RegisterRoutes(api4debtus.RegisterHttpRoutes),
		module.RegisterDelays(facade4debtus.InitDelays4debtus),
	)
}
