package spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/api4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/const4spaceus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4spaceus.ModuleID, module.RegisterRoutes(api4spaceus.RegisterHttpRoutes))
}
