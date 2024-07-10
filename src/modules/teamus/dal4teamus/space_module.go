package dal4teamus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/core/coremodels"
)

const SpaceModulesCollection = coremodels.ModulesCollection

func NewSpaceModuleKey(teamID, moduleID string) *dal.Key {
	teamKey := NewSpaceKey(teamID)
	return dal.NewKeyWithParentAndID(teamKey, SpaceModulesCollection, moduleID)
}

func NewSpaceModuleEntry[D any](teamID, moduleID string, data D) record.DataWithID[string, D] {
	key := NewSpaceModuleKey(teamID, moduleID)
	return record.NewDataWithID(moduleID, key, data)
}
