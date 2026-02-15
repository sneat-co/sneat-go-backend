package calendarium

import (
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/api4calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/const4calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dbo4calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/delays4calendarium"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	facade4linkage.RegisterDboFactory(const4calendarium.ExtensionID, const4calendarium.HappeningsCollection,
		facade4linkage.NewDboFactory(
			func() facade4linkage.SpaceItemDboWithRelatedAndIDs {
				return new(dbo4calendarium.HappeningDbo)
			},
			func() dal4spaceus.SpaceModuleDbo {
				return new(dbo4calendarium.CalendariumSpaceDbo)
			},
		),
	)
	return extension.NewExtension(const4calendarium.ExtensionID,
		extension.RegisterRoutes(api4calendarium.RegisterHttpRoutes),
		extension.RegisterDelays(delays4calendarium.InitDelaying),
	)
}
