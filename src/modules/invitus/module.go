package invitus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/api4invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/const4invitus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4invitus.ModuleID, module.RegisterRoutes(api4invitus.RegisterHttpRoutes))
}
