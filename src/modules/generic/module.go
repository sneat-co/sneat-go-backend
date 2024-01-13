package generic

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/generic/api4generic"
	"github.com/sneat-co/sneat-go-backend/src/modules/generic/const4generic"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4generic.ModuleID, api4generic.RegisterHttpRoutes)
}
