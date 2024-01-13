package invitus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/api4invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/const4invitus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4invitus.ModuleID, api4invitus.RegisterHttpRoutes)
}
