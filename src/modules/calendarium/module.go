package calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/api4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/delays4calendarium"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4calendarium.ModuleID,
		module.RegisterRoutes(api4calendarium.RegisterHttpRoutes),
		module.RegisterDelays(delays4calendarium.InitDelaying),
	)
}
