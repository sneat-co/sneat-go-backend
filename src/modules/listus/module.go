package listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/api4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4listus.ModuleID, api4listus.RegisterHttpRoutes)
}
