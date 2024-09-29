package linkage

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/api4linkage"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/linkage/const4linkage"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4linkage.ModuleID, module.RegisterRoutes(api4linkage.RegisterHttpRoutes))
}
