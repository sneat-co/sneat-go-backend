package contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/api4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/const4contactus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/contactus/delays4contactus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4contactus.ModuleID,
		module.RegisterRoutes(api4contactus.RegisterHttpRoutes),
		module.RegisterDelays(delays4contactus.InitDelays4contactus),
	)
}
