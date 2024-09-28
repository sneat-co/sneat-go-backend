package facade4auth

import (
	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/models4auth"
)

//func newUserGooglePlusKey(id string) *dal.Key {
//	return dal.NewKeyWithID(models.UserGooglePlusKind, id)
//}

type UserGooglePlusDalGae struct {
}

func NewUserGooglePlusDalGae() UserGooglePlusDalGae {
	return UserGooglePlusDalGae{}
}

func (UserGooglePlusDalGae) GetUserGooglePlusByID(_ context.Context, id string) (userGooglePlus models4auth.UserGooglePlus, err error) {
	//var userGooglePlusEntity models.UserGooglePlusEntity
	//if err = gaedb.Get(ctx, newUserGooglePlusKey(id), &userGooglePlusEntity); err != nil {
	//	if err == datastore.ErrNoSuchEntity {
	//		err = db.NewErrNotFoundByStrID(models.UserGooglePlusKind, id, err)
	//	}
	//	return
	//}
	//userGooglePlus = models.UserGooglePlus{StringID: db.StringID{ContactID: id}, UserGooglePlusEntity: &userGooglePlusEntity}
	err = errors.New("not implemented")
	return
}
