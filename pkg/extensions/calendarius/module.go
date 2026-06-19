package calendarius

import (
	"github.com/sneat-co/sneat-core-modules/linkage/facade4linkage"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/api4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/const4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dbo4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/delays4calendarius"
	"github.com/sneat-co/sneat-go-core/extension"
)

func Extension() extension.Config {
	facade4linkage.RegisterDboFactory(const4calendarius.ExtensionID, const4calendarius.HappeningsCollection,
		facade4linkage.NewDboFactory(
			func() facade4linkage.SpaceItemDboWithRelatedAndIDs {
				return new(dbo4calendarius.HappeningDbo)
			},
			func() dal4spaceus.SpaceModuleDbo {
				return new(dbo4calendarius.CalendariusSpaceDbo)
			},
		),
	)
	return extension.NewExtension(const4calendarius.ExtensionID,
		extension.RegisterRoutes(api4calendarius.RegisterHttpRoutes),
		extension.RegisterDelays(delays4calendarius.InitDelaying),
	)
}
