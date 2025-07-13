package deedus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/deedus/api4deedus"
	"github.com/sneat-co/sneat-go-backend/src/modules/deedus/const4deedus"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	return module.NewExtension(const4deedus.ModuleID, module.RegisterRoutes(api4deedus.RegisterHttpRoutes))
}
