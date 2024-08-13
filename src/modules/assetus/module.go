package assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/api4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4assetus.ModuleID, module.RegisterRoutes(api4assetus.RegisterHttpRoutes))
}
