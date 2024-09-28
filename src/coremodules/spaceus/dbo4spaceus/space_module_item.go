package dbo4spaceus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
)

func NewSpaceModuleItemKey[K comparable](spaceID, moduleID, collection string, itemID K) *dal.Key {
	teamModuleKey := NewSpaceModuleKey(spaceID, moduleID)
	return dal.NewKeyWithParentAndID(teamModuleKey, collection, itemID)
}

func NewSpaceModuleItemKeyFromItemRef(itemRef dbo4linkage.SpaceModuleItemRef) *dal.Key {
	return NewSpaceModuleItemKey(itemRef.Space, itemRef.Module, itemRef.Collection, itemRef.ItemID)
}
