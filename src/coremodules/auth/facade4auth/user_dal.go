package facade4auth

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/random"
	"strings"
	"time"
)

type UserDalGae struct {
}

func (UserDalGae) DelaySetUserPreferredLocale(_ context.Context, delay time.Duration, userID string, localeCode5 string) error {
	//return delaySetUserPreferredLocale.EnqueueWork(ctx, delaying.With(common4debtus.QUEUE_USERS, "set-user-preferred-locale", delay), userID, localeCode5)
	//TODO implement me
	panic("implement me")
}

func NewUserDalGae() UserDalGae {
	return UserDalGae{}
}

var _ UserDal = (*UserDalGae)(nil)

func (userDal UserDalGae) GetUserByStrID(ctx context.Context, userID string) (user dbo4userus.UserEntry, err error) {
	user = dbo4userus.NewUserEntry(userID)
	return user, dal4userus.GetUser(ctx, nil, user)
}

func (userDal UserDalGae) GetUserByVkUserID(_ context.Context, vkUserID int64) (dbo4userus.UserEntry, error) {
	panic("not implemented")
	//query := datastore.NewQuery(models.AppUserKind).Filter("VkUserID =", vkUserID)
	//return userDal.getUserByQuery(ctx, query, "VkUserID")
}

//func (userDal UserDalGae) getUserByQuery(ctx context.Context, query dal.Query, searchCriteria string) (appUser models4debtus.AppUserOBSOLETE, err error) {
//	userEntities := make([]*models4debtus.DebutsAppUserDataOBSOLETE, 0, 2)
//	var db dal.DB
//	if db, err = facade.GetSneatDB(c); err != nil {
//		return
//	}
//	var userRecords []dal.Record
//
//	if userRecords, err = db.QueryAllRecords(ctx, query); err != nil {
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

func (userDal UserDalGae) CreateAnonymousUser(ctx context.Context) (user dbo4userus.UserEntry, err error) {
	return userDal.CreateUser(ctx, &dbo4userus.UserDbo{
		IsAnonymous: true,
	})
}

func (userDal UserDalGae) CreateUser(ctx context.Context, userData *dbo4userus.UserDbo) (user dbo4userus.UserEntry, err error) {
	user = dbo4userus.NewUserEntryWithDbo("", userData)

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) error {
		if err = tx.Insert(ctx, user.Record); err != nil {
			return err
		}
		user.ID = user.Record.Key().ID.(string)
		user.Data = user.Record.Data().(*dbo4userus.UserDbo)
		return nil
	})
	return
}

func GenerateRandomUserID(ctx context.Context, tx dal.ReadwriteTransaction) (userID string, err error) {

	const maxAttempts = 5

	randomIDs := make([]string, 0, maxAttempts)

	for i := 1; i <= maxAttempts; i++ {
		userID = random.ID(12)
		userKey := dbo4userus.NewUserKey(userID)
		userData := make(map[string]any)
		userRecord := dal.NewRecordWithData(userKey, userData)
		if err = tx.Get(ctx, userRecord); err != nil {
			if dal.IsNotFound(err) {
				err = nil
				return
			}
			return "", fmt.Errorf("failed to check user record exists for a random userID: %w", err)
		}
	}
	return "", fmt.Errorf("too many attempts (%d) to generate a random userID, tried next IDs: %s", maxAttempts, strings.Join(randomIDs, ", "))
}
