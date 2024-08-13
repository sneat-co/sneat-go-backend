package userus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/api4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/delays4userus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewModule(const4userus.ModuleID,
		module.RegisterRoutes(api4userus.RegisterHttpRoutes),
		module.RegisterDelays(delays4userus.InitDelays4userus),
	)
}
