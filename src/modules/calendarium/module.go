package calendarium

import (
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/api4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/delays4calendarium"
	"github.com/sneat-co/sneat-go-core/module"
)

func Module() module.Module {
	facade4linkage.RegisterDboFactory(const4calendarium.ModuleID, const4calendarium.HappeningsCollection,
		facade4linkage.NewDboFactory(
			func() facade4linkage.SpaceItemDboWithRelatedAndIDs {
				return new(dbo4calendarium.HappeningDbo)
			},
			func() dal4spaceus.SpaceModuleDbo {
				return new(dbo4calendarium.CalendariumSpaceDbo)
			},
		),
	)
	return module.NewExtension(const4calendarium.ModuleID,
		module.RegisterRoutes(api4calendarium.RegisterHttpRoutes),
		module.RegisterDelays(delays4calendarium.InitDelaying),
	)
}
