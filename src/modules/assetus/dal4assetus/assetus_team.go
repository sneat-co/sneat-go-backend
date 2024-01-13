package dal4assetus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/const4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/models4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
)

type AssetusTeamContext = record.DataWithID[string, *models4assetus.AssetusTeamDto]

// AssetsCollection is a name of a collection in DB
const AssetsCollection = "assets"

var AssetusRootKey = dal.NewKeyWithID(dal4teamus.TeamModulesCollection, const4assetus.ModuleID)
