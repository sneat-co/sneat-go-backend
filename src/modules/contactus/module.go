package contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/api4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-core/modules"
)

func Module() modules.Module {
	return modules.NewModule(const4contactus.ModuleID, api4contactus.RegisterHttpRoutes)
}
