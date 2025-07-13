package listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/api4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/const4listus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewExtension(const4listus.ModuleID, module.RegisterRoutes(api4listus.RegisterHttpRoutes))
}
