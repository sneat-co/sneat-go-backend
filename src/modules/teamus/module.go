package teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/api4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/const4teamus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4teamus.ModuleID, api4teamus.RegisterHttpRoutes)
}
