package facade4auth

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/auth/models4auth"
	"strconv"
)

func NewUserVkKey(vkUserID int64) *dal.Key {
	return dal.NewKeyWithID(models4auth.UserVkKind, strconv.FormatInt(vkUserID, 10))
}

type UserVkDalGae struct {
}

func NewUserVkDalGae() UserVkDalGae {
	return UserVkDalGae{}
}

func (UserVkDalGae) GetUserVkByID(c context.Context, vkUserID int64) (vkUser models4auth.UserVk, err error) {
	//vkUserKey := NewUserVkKey(c, vkUserID)
	//var vkUserEntity models.UserVkEntity
	//if err = gaedb.Get(c, vkUserKey, &vkUserEntity); err != nil {
	//	if err == datastore.ErrNoSuchEntity {
	//		err = db.NewErrNotFoundByIntID(models.UserVkKind, vkUserID, nil)
	//	}
	//	return
	//}
	//vkUser = models.UserVk{IntegerID: db.NewIntID(vkUserID), UserVkEntity: &vkUserEntity}
	return vkUser, errors.New("not implemented")
}

func (UserVkDalGae) SaveUserVk(c context.Context, userVk models4auth.UserVk) (err error) {
	//k := NewUserVkKey(c, userVk.ContactID)
	//_, err = gaedb.Put(c, k, userVk)
	return errors.New("not implemented")
}
