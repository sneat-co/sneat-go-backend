package spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/api4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/const4spaceus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4spaceus.ModuleID, api4spaceus.RegisterHttpRoutes)
}
