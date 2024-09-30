package modules

import (
	"github.com/sneat-co/sneat-core-modules"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/deedus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus"
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus"
	"github.com/sneat-co/sneat-go-backend/src/modules/sportus"
	"github.com/sneat-co/sneat-go-core/module"
)

func standardModules() []module.Module {
	return []module.Module{
		assetus.Module(),
		deedus.Module(),
		calendarium.Module(),
		listus.Module(),
		retrospectus.Module(),
		scrumus.Module(),
		sportus.Module(),
	}
}

func Modules() []module.Module {
	return append(
		sneat_core_modules.CoreModules(),
		standardModules()...,
	)
}
