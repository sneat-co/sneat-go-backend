package facade4retrospectus

import (
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
)

func getUserRetroRecord(uid, team string, data interface{}) dal.Record {
	userKey := dbo4userus.NewUserKey(uid)
	key := dal.NewKeyWithParentAndID(userKey, "api4meetingus", team)
	return dal.NewRecordWithData(key, data)
}
