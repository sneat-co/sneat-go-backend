package retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/api4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/const4retrospectus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4retrospectus.ModuleID, api4retrospectus.RegisterHttpRoutes)
}
