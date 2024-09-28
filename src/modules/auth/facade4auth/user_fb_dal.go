package facade4auth

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/models4auth"
)

func NewUserFacebookKey(ctx context.Context, fbAppOrPageID, fbUserOrPageScopeID string) *dal.Key {
	if fbAppOrPageID == "" {
		panic("fbAppOrPageID is empty string")
	}
	if fbUserOrPageScopeID == "" {
		panic("fbUserOrPageScopeID is empty string")
	}
	return dal.NewKeyWithID(models4auth.UserFacebookCollection, fbAppOrPageID+":"+fbUserOrPageScopeID)
}

type UserFacebookDalGae struct {
}

func NewUserFacebookDalGae() UserFacebookDalGae {
	return UserFacebookDalGae{}
}

func (UserFacebookDalGae) SaveFbUser(_ context.Context, tx dal.ReadwriteTransaction, fbUser models4auth.UserFacebook) (err error) {
	//key := NewUserFacebookKey(ctx, fbUser.FbAppOrPageID, fbUser.FbUserOrPageScopeID)
	//if _, err = gaedb.Put(c, key, fbUser.Data); err != nil {
	//	return
	//}
	//return
	return errors.New("not implemented")
}

func (UserFacebookDalGae) DeleteFbUser(_ context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (err error) {
	//key := NewUserFacebookKey(ctx, fbAppOrPageID, fbUserOrPageScopeID)
	//if err = gaedb.Delete(ctx, key); err != nil {
	//	return
	//}
	return errors.New("not implemented")
}

func (UserFacebookDalGae) GetFbUserByFbID(_ context.Context, fbAppOrPageID, fbUserOrPageScopeID string) (fbUser models4auth.UserFacebook, err error) {
	err = errors.New("not implemented")
	return
	//var entity models.UserFacebookData
	//if err = gaedb.Get(ctx, NewUserFacebookKey(ctx, fbAppOrPageID, fbUserOrPageScopeID), &entity); err != nil {
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
