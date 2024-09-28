package dbo4spaceus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/core/coremodels"
)

const SpaceModulesCollection = coremodels.ModulesCollection

func NewSpaceModuleKey(spaceID, moduleID string) *dal.Key {
	teamKey := NewSpaceKey(spaceID)
	return dal.NewKeyWithParentAndID(teamKey, SpaceModulesCollection, moduleID)
}

func NewSpaceModuleEntry[D any](spaceID, moduleID string, data D) record.DataWithID[string, D] {
	key := NewSpaceModuleKey(spaceID, moduleID)
	return record.NewDataWithID(moduleID, key, data)
}
