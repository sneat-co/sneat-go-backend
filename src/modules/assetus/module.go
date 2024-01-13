package assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/api4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4assetus.ModuleID, api4assetus.RegisterHttpRoutes)
}
