package scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/api4scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/const4srumus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4srumus.ModuleID, api4scrumus.RegisterHttpRoutes)
}
