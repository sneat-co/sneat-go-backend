package dal4assetus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-core-modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dbo4assetus"
)

type AssetusSpaceEntry = record.DataWithID[string, *dbo4assetus.AssetusSpaceDbo]

// AssetsCollection is a name of a collection in DB
const AssetsCollection = "assets"

var AssetusRootKey = dal.NewKeyWithID(dbo4spaceus.SpaceModulesCollection, const4assetus.ExtensionID)
