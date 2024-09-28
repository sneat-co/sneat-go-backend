package dal4userus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	"github.com/sneat-co/sneat-go-backend/src/core/coremodels"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
)

const UserModulesCollection = coremodels.ModulesCollection

func NewUserModuleKey(userID, moduleID string) *dal.Key {
	userKey := dbo4userus.NewUserKey(userID)
	return dal.NewKeyWithParentAndID(userKey, UserModulesCollection, moduleID)
}

func NewUserModuleEntry[D any](userID, moduleID string, data D) record.DataWithID[string, D] {
	key := NewUserModuleKey(userID, moduleID)
	return record.NewDataWithID(moduleID, key, data)
}
