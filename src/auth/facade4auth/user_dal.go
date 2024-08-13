package facade4auth

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"time"
)

type UserDalGae struct {
}

func (UserDalGae) DelaySetUserPreferredLocale(c context.Context, delay time.Duration, userID string, localeCode5 string) error {
	//return delaySetUserPreferredLocale.EnqueueWork(c, delaying.With(common4debtus.QUEUE_USERS, "set-user-preferred-locale", delay), userID, localeCode5)
	//TODO implement me
	panic("implement me")
}

func NewUserDalGae() UserDalGae {
	return UserDalGae{}
}

var _ UserDal = (*UserDalGae)(nil)

func (userDal UserDalGae) GetUserByStrID(c context.Context, userID string) (user dbo4userus.UserEntry, err error) {
	user = dbo4userus.NewUserEntry(userID)
	return user, dal4userus.GetUser(c, nil, user)
}

func (userDal UserDalGae) GetUserByVkUserID(c context.Context, vkUserID int64) (dbo4userus.UserEntry, error) {
	panic("not implemented")
	//query := datastore.NewQuery(models.AppUserKind).Filter("VkUserID =", vkUserID)
	//return userDal.getUserByQuery(c, query, "VkUserID")
}

//func (userDal UserDalGae) getUserByQuery(c context.Context, query dal.Query, searchCriteria string) (appUser models4debtus.AppUserOBSOLETE, err error) {
//	userEntities := make([]*models4debtus.DebutsAppUserDataOBSOLETE, 0, 2)
//	var db dal.DB
//	if db, err = facade.GetDatabase(c); err != nil {
//		return
//	}
//	var userRecords []dal.Record
//
//	if userRecords, err = db.QueryAllRecords(c, query); err != nil {
//		return
//	}
//	switch len(userRecords) {
//	case 1:
//		logus.Debugf(c, "getUserByQuery(%v) => %v: %v", searchCriteria, userRecords[0].Key().ContactID, userEntities[0])
//		ur := userRecords[0]
//		return models4debtus.NewAppUserOBSOLETE(ur.Key().ContactID.(string), ur.Data().(*models4debtus.DebutsAppUserDataOBSOLETE)), nil
//	case 0:
//		err = dal.ErrRecordNotFound
//		logus.Debugf(c, "getUserByQuery(%v) => %v", searchCriteria, err)
//		return
//	default: // > 1
//		errDup := dal.ErrDuplicateUser{ // TODO: ErrDuplicateUser should be moved out from dalgo
//			SearchCriteria:   searchCriteria,
//			DuplicateUserIDs: make([]string, len(userRecords)),
//		}
//		for i, userRecord := range userRecords {
//			errDup.DuplicateUserIDs[i] = userRecord.Key().ContactID.(string)
//		}
//		err = errDup
//		return
//	}
//}

func (userDal UserDalGae) CreateAnonymousUser(c context.Context) (user dbo4userus.UserEntry, err error) {
	return userDal.CreateUser(c, &dbo4userus.UserDbo{
		IsAnonymous: true,
	})
}

func (userDal UserDalGae) CreateUser(c context.Context, userData *dbo4userus.UserDbo) (user dbo4userus.UserEntry, err error) {
	user = dbo4userus.NewUserEntryWithDbo("", userData)

	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if err = tx.Insert(c, user.Record); err != nil {
			return err
		}
		user.ID = user.Record.Key().ID.(string)
		user.Data = user.Record.Data().(*dbo4userus.UserDbo)
		return nil
	})
	return
}
