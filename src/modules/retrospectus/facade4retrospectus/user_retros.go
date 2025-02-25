package facade4retrospectus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/coretypes"
)

func getUserRetroRecord(uid string, spaceID coretypes.SpaceID, data interface{}) dal.Record {
	userKey := dbo4userus.NewUserKey(uid)
	key := dal.NewKeyWithParentAndID(userKey, "api4meetingus", spaceID)
	return dal.NewRecordWithData(key, data)
}
