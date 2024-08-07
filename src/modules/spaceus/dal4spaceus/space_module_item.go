package dal4spaceus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/linkage/dbo4linkage"
	"github.com/strongo/random"
)

func NewSpaceModuleItemKey[K comparable](teamID, moduleID, collection string, itemID K) *dal.Key {
	teamModuleKey := NewSpaceModuleKey(teamID, moduleID)
	return dal.NewKeyWithParentAndID(teamModuleKey, collection, itemID)
}

func NewSpaceModuleItemKeyFromItemRef(itemRef dbo4linkage.SpaceModuleItemRef) *dal.Key {
	return NewSpaceModuleItemKey(itemRef.Space, itemRef.Module, itemRef.Collection, itemRef.ItemID)
}

func GenerateNewSpaceModuleItemKey(ctx context.Context, tx dal.ReadwriteTransaction, teamID, moduleID, collection string, length, maxAttempts int) (id string, key *dal.Key, err error) {
	for i := 0; i < maxAttempts; i++ {
		id = random.ID(length)
		key = NewSpaceModuleItemKey(teamID, moduleID, collection, id)
		record := dal.NewRecordWithData(key, make(map[string]interface{}))
		if err := tx.Get(ctx, record); err != nil { // TODO: use tx.Exists()
			if dal.IsNotFound(err) {
				return id, key, nil
			}
			return "", nil, err
		}
	}
	return "", nil, errors.New("too many attempts  to generate a random happening ContactID")
}
