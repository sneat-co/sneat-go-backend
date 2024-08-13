package scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/api4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/const4srumus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4srumus.ModuleID, module.RegisterRoutes(api4scrumus.RegisterHttpRoutes))
}
