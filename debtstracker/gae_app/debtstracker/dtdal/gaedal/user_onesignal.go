package gaedal

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type UserOneSignalDalGae struct {
}

func NewUserOneSignalDalGae() UserOneSignalDalGae {
	return UserOneSignalDalGae{}
}

func (userOneSignalDalGae UserOneSignalDalGae) SaveUserOneSignal(c context.Context, userID string, oneSignalUserID string) (userOneSignal models.UserOneSignal, err error) {
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
	//userOneSignal = models.UserOneSignal{StringID: db.StringID{ID: oneSignalUserID}, UserOneSignalEntity: &entity}
	return userOneSignal, errors.New("not implemented")
}

func (UserOneSignalDalGae) NewUserOneSignalKey(oneSignalUserID string) *dal.Key {
	return dal.NewKeyWithID(models.UserOneSignalKind, oneSignalUserID)
}
