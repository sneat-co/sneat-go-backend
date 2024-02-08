package gaedal

import (
	"context"
	"github.com/crediterra/money"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/delaying"
	"github.com/strongo/log"
	"strings"
	"time"
)

type UserDalGae struct {
}

func NewUserDalGae() UserDalGae {
	return UserDalGae{}
}

var _ dtdal.UserDal = (*UserDalGae)(nil)

func (userDal UserDalGae) SetLastCurrency(c context.Context, userID string, currency money.CurrencyCode) (err error) {
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return err
	}
	return db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		user, err := facade.User.GetUserByID(c, tx, userID)
		if err != nil {
			return err
		}
		user.Data.SetLastCurrency(string(currency))
		return facade.User.SaveUser(c, tx, user)
	})
}

func (userDal UserDalGae) GetUserByStrID(c context.Context, userID string) (user models.AppUser, err error) {
	return facade.User.GetUserByID(c, nil, userID)
}

func (userDal UserDalGae) GetUserByVkUserID(c context.Context, vkUserID int64) (models.AppUser, error) {
	panic("not implemented")
	//query := datastore.NewQuery(models.AppUserKind).Filter("VkUserID =", vkUserID)
	//return userDal.getUserByQuery(c, query, "VkUserID")
}

func (userDal UserDalGae) GetUserByEmail(c context.Context, email string) (models.AppUser, error) {
	email = strings.ToLower(email)
	query := dal.From(models.AppUserKind).Where(
		dal.WhereField("EmailAddress", dal.Equal, email),
		dal.WhereField("EmailConfirmed", dal.Equal, true),
	).Limit(2).SelectInto(models.NewAppUserRecord)
	user, err := userDal.getUserByQuery(c, query, "EmailAddress, is confirmed")
	if user.ID == "" && dal.IsNotFound(err) {
		query = dal.From(models.AppUserKind).
			Where(
				dal.WhereField("EmailAddress", dal.Equal, email),
				dal.WhereField("EmailConfirmed", dal.Equal, false),
			).
			Limit(2).
			SelectInto(models.NewAppUserRecord)
		user, err = userDal.getUserByQuery(c, query, "EmailAddress, is not confirmed")
	}
	log.Debugf(c, "GetUserByEmail() => err=%v, User(id=%d): %v", err, user.ID, user)
	return user, err
}

func (userDal UserDalGae) getUserByQuery(c context.Context, query dal.Query, searchCriteria string) (appUser models.AppUser, err error) {
	userEntities := make([]*models.DebutsAppUserDataOBSOLETE, 0, 2)
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	var userRecords []dal.Record

	if userRecords, err = db.QueryAllRecords(c, query); err != nil {
		return
	}
	switch len(userRecords) {
	case 1:
		log.Debugf(c, "getUserByQuery(%v) => %v: %v", searchCriteria, userRecords[0].Key().ID, userEntities[0])
		ur := userRecords[0]
		return models.NewAppUser(ur.Key().ID.(string), ur.Data().(*models.DebutsAppUserDataOBSOLETE)), nil
	case 0:
		err = dal.ErrRecordNotFound
		log.Debugf(c, "getUserByQuery(%v) => %v", searchCriteria, err)
		return
	default: // > 1
		errDup := dal.ErrDuplicateUser{ // TODO: ErrDuplicateUser should be moved out from dalgo
			SearchCriteria:   searchCriteria,
			DuplicateUserIDs: make([]string, len(userRecords)),
		}
		for i, userRecord := range userRecords {
			errDup.DuplicateUserIDs[i] = userRecord.Key().ID.(string)
		}
		err = errDup
		return
	}
}

func (userDal UserDalGae) CreateAnonymousUser(c context.Context) (user models.AppUser, err error) {
	return userDal.CreateUser(c, &models.DebutsAppUserDataOBSOLETE{
		IsAnonymous: true,
	})
}

func (userDal UserDalGae) CreateUser(c context.Context, userData *models.DebutsAppUserDataOBSOLETE) (user models.AppUser, err error) {
	user = models.NewAppUser("", userData)

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if err = tx.Insert(c, user.Record); err != nil {
			return err
		}
		user.ID = user.Record.Key().ID.(string)
		user.Data = user.Record.Data().(*models.DebutsAppUserDataOBSOLETE)
		return nil
	})
	return
}

func (UserDalGae) DelayUpdateUserWithBill(c context.Context, userID string, billID string) (err error) {
	if err = delayUpdateUserWithBill.EnqueueWork(c, delaying.With(common.QUEUE_BILLS, "UpdateUserWithBill", 0), userID, billID); err != nil {
		return
	}
	return
}

func (UserDalGae) DelayUpdateUserWithContact(c context.Context, userID string, contactID string) (err error) {
	if err = delayedUpdateUserWithContact.EnqueueWork(c, delaying.With(common.QUEUE_USERS, "updateUserWithContact", time.Second/10), userID, contactID); err != nil {
		return
	}
	return
}

func updateUserWithContact(c context.Context, userID, contactID string) (err error) {
	log.Debugf(c, "updateUserWithContact(userID=%v, contactID=%v)", userID, contactID)
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		return
	}
	return db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		var contact models.Contact
		if contact, err = facade.GetContactByID(c, tx, contactID); err != nil {
			if dal.IsNotFound(err) {
				log.Warningf(c, "contact not found: %v", err)
				return nil
			}
			log.Errorf(c, "updateUserWithContact: %v", err)
			return
		}
		var user models.AppUser

		if user, err = facade.User.GetUserByID(c, tx, userID); err != nil {
			return
		}
		if dal.IsNotFound(err) {
			log.Errorf(c, err.Error())
			err = nil
		}

		if _, changed := models.AddOrUpdateContact(&user, contact); changed {
			if err = facade.User.SaveUser(c, tx, user); err != nil {
				return
			}
		} else {
			log.Debugf(c, "user not changed")
		}
		return
	})
}
