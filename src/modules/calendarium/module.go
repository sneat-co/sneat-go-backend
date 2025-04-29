package calendarium

import (
	"github.com/sneat-co/sneat-core-modules/linkage/dbo4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-core-modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/api4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/delays4calendarium"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	facade4spaceus.RegisterDboFactory(const4calendarium.ModuleID, const4calendarium.HappeningsCollection,
		func() (dal4spaceus.SpaceModuleDbo, dal4spaceus.SpaceItemDbo, *dbo4linkage.WithRelatedAndIDs) {
			dbo := new(dbo4calendarium.HappeningDbo)
			return new(dbo4calendarium.CalendariumSpaceDbo), dbo, nil
		},
	)
	return module.NewModule(const4calendarium.ModuleID,
		module.RegisterRoutes(api4calendarium.RegisterHttpRoutes),
		module.RegisterDelays(delays4calendarium.InitDelaying),
	)
}
