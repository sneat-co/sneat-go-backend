package retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/api4retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/const4retrospectus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewExtension(const4retrospectus.ModuleID, module.RegisterRoutes(api4retrospectus.RegisterHttpRoutes))
}
