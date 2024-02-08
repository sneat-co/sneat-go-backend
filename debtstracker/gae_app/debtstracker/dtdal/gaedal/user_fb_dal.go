package gaedal

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func NewUserFacebookKey(c context.Context, fbAppOrPageID, fbUserOrPageScopeID string) *dal.Key {
	if fbAppOrPageID == "" {
		panic("fbAppOrPageID is empty string")
	}
	if fbUserOrPageScopeID == "" {
		panic("fbUserOrPageScopeID is empty string")
	}
	return dal.NewKeyWithID(models.UserFacebookCollection, fbAppOrPageID+":"+fbUserOrPageScopeID)
}

type UserFacebookDalGae struct {
}

func NewUserFacebookDalGae() UserFacebookDalGae {
	return UserFacebookDalGae{}
}

func (UserFacebookDalGae) SaveFbUser(c context.Context, tx dal.ReadwriteTransaction, fbUser models.UserFacebook) (err error) {
	//key := NewUserFacebookKey(c, fbUser.FbAppOrPageID, fbUser.FbUserOrPageScopeID)
	//if _, err = gaedb.Put(c, key, fbUser.Data); err != nil {
	//	return
	//}
	//return
	return errors.New("not implemented")
}

func (UserFacebookDalGae) DeleteFbUser(c context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (err error) {
	//key := NewUserFacebookKey(c, fbAppOrPageID, fbUserOrPageScopeID)
	//if err = gaedb.Delete(c, key); err != nil {
	//	return
	//}
	return errors.New("not implemented")
}

func (UserFacebookDalGae) GetFbUserByFbID(c context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (fbUser models.UserFacebook, err error) {
	err = errors.New("not implemented")
	return
	//var entity models.UserFacebookData
	//if err = gaedb.Get(c, NewUserFacebookKey(c, fbAppOrPageID, fbUserOrPageScopeID), &entity); err != nil {
	//	if err == datastore.ErrNoSuchEntity {
	//		err = db.NewErrNotFoundByStrID(models.UserFacebookCollection, fbUserOrPageScopeID, err)
	//	}
	//	return
	//}
	//fbUser = models.UserFacebook{
	//	FbAppOrPageID:       fbAppOrPageID,
	//	FbUserOrPageScopeID: fbUserOrPageScopeID,
	//	UserFacebookEntity:  &entity,
	//}
	//return
}
