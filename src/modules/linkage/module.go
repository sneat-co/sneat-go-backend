package linkage

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/api4linkage"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/const4linkage"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4linkage.ModuleID, api4linkage.RegisterHttpRoutes)
}
