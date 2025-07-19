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
	"github.com/sneat-co/sneat-go-core/extension"
)

func standardModules() []extension.Config {
	return []extension.Config{
		assetus.Extension(),
		deedus.Extension(),
		calendarium.Extension(),
		listus.Extension(),
		retrospectus.Extension(),
		scrumus.Extension(),
		sportus.Extension(),
	}
}

func Extensions() []extension.Config {
	return append(
		sneat_core_modules.CoreExtensions(),
		standardModules()...,
	)
}
