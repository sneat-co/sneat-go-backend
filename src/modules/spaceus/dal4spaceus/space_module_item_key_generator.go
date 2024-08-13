package dal4spaceus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/strongo/random"
)

func GenerateNewSpaceModuleItemKey(ctx context.Context, tx dal.ReadwriteTransaction, teamID, moduleID, collection string, length, maxAttempts int) (id string, key *dal.Key, err error) {
	for i := 0; i < maxAttempts; i++ {
		id = random.ID(length)
		key = dbo4spaceus.NewSpaceModuleItemKey(teamID, moduleID, collection, id)
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
