package facade4userus

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
)

func SaveUserOneSignal(c context.Context, userID string, oneSignalUserID string) (userOneSignal models4auth.UserOneSignal, err error) {
	//key := userOneSignalDalGae.NewUserOneSignalKey(c, oneSignalUserID)
	//var entity models.UserOneSignalEntity
	//// Save if no entity or AppUserIntID changed
	//if err = gaedb.Get(c, key, &entity); err == datastore.ErrNoSuchEntity || entity.UserID != userID {
	//	entity = models.UserOneSignalEntity{UserID: userID, Created: time.Now()}
	//	if _, err = gaedb.Put(c, key, &entity); err != nil {
	//		return
	//	}
	//} else if err != nil {
	//	return
	//}
	//userOneSignal = models.UserOneSignal{StringID: db.StringID{ContactID: oneSignalUserID}, UserOneSignalEntity: &entity}
	return userOneSignal, errors.New("not implemented")
}

func NewUserOneSignalKey(oneSignalUserID string) *dal.Key {
	return dal.NewKeyWithID(models4auth.UserOneSignalKind, oneSignalUserID)
}