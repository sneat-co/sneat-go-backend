package facade4debtus

import (
	"context"
	"errors"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/dalgo/record"
	models4auth2 "github.com/sneat-co/sneat-core-modules/auth/models4auth"
	"github.com/sneat-co/sneat-core-modules/auth/unsorted4auth"
	"github.com/sneat-co/sneat-core-modules/common4all"
	"github.com/sneat-co/sneat-core-modules/core/sneaterrors"
	"github.com/sneat-co/sneat-core-modules/userus/dal4userus"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"strings"
	"time"
)

type userFacade struct {
}

var User = userFacade{}

// Deprecated: use facade4userus instead
func (userFacade) GetUserByIdOBSOLETE(ctx context.Context, tx dal.ReadSession, userID string) (user models4debtus.AppUserOBSOLETE, err error) {
	if tx == nil {
		if tx, err = facade.GetSneatDB(ctx); err != nil {
			return
		}
	}

	key := dal.NewKeyWithID(models4debtus.AppUserKind, userID)
	user.Data = new(models4debtus.DebutsAppUserDataOBSOLETE)
	user.WithID = record.WithID[string]{
		ID:     userID,
		Key:    key,
		Record: dal.NewRecordWithData(key, user.Data),
	}
	err = tx.Get(ctx, user.Record)
	return
}

func (userFacade) GetUsersByIDs(ctx context.Context, userIDs []string) (users []dbo4userus.UserEntry, err error) {
	return dal4userus.GetUsersByIDs(ctx, userIDs)
}

func (uf userFacade) CreateUserByEmail(
	ctx context.Context,
	email, name string,
) (
	user dbo4userus.UserEntry,
	userEmail models4auth2.UserEmailEntry,
	err error,
) {
	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if userEmail, err = unsorted4auth.UserEmail.GetUserEmailByID(ctx, tx, email); err == nil {
			return sneaterrors.ErrEmailAlreadyRegistered
		} else if !dal.IsNotFound(err) {
			return
		}

		if userEmail.ID == "" {
			logus.Errorf(ctx, "userEmail.ContactID is empty string")
			userEmail.ID = strings.ToLower(strings.TrimSpace(email))
		}

		//userData := unsorted4auth.CreateUserEntity(unsorted4auth.CreateUserData{
		//	ScreenName: name,
		//})
		//userData.AddAccount(userEmail.UserAccount())

		err = errors.New("not implemented")
		return
		//if user, err = facade4auth.UserEntry.CreateUser(ctx, userData); err != nil {
		//	return
		//}

		//userEmail.Data.Provider = "email"
		//userEmail.Data.EmailLowerCase = email
		//if err = userEmail.Data.SetPassword(dtdal.RandomCode(8)); err != nil {
		//	return
		//}
		//
		//err = facade4auth.UserEmail.SaveUserEmail(ctx, tx, userEmail)
		//return
	})

	return
}

// GetOrCreateEmailUser is used in invites.
func (uf userFacade) GetOrCreateEmailUser(
	ctx context.Context,
	email string,
	isConfirmed bool,
	createUserData *unsorted4auth.CreateUserData,
	clientInfo common4all.ClientInfo,
) (
	userEmail models4auth2.UserEmailEntry,
	isNewUser bool,
	err error,
) {

	var appUser dbo4userus.UserEntry

	err = facade.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
		if userEmail, err = unsorted4auth.UserEmail.GetUserEmailByID(ctx, tx, email); err == nil {
			return // UserEntry found
		} else if !dal.IsNotFound(err) { //
			return // Internal error
		}
		//err = nil // Clear dtdal.ErrRecordNotFound

		now := time.Now()
		isNewUser = true
		userEmail = models4auth2.NewUserEmail(email, models4auth2.NewUserEmailData(0, isConfirmed, "email"))
		appUser = dbo4userus.NewUserEntry("") //clientInfo
		appUser.Data.CreatedAt = now
		appUser.Data.AddAccount(userEmail.UserAccount())

		//var to db.RunOptions = dtdal.CrossGroupTransaction

		if err = tx.Set(ctx, appUser.Record); err != nil {
			return fmt.Errorf("failed to save new appUser to datastore: %w", err)
		}
		userEmail.Data.CreatedAt = now

		if err = unsorted4auth.UserEmail.SaveUserEmail(ctx, tx, userEmail); err != nil {
			return err
		}
		return nil
	})

	return
}

//func (uf userFacade) GetOrCreateUserGoogleOnSignIn(
//	ctx context.Context, googleUser *gae_user.UserEntry, appUserID string, clientInfo models.ClientInfo,
//) (
//	userGoogle models.UserAccountEntry, appUser models.AppUserOBSOLETE, err error,
//) {
//	if googleUser == nil {
//		panic("googleUser == nil")
//	}
//	getUserAccountRecordFromDB := func(ctx context.Context) (appuser.AccountRecord, error) {
//		userGoogle, err = dtdal.UserGoogle.GetUserGoogleByID(ctx, googleUser.ContactID)
//		return &userGoogle, err
//	}
//	newUserAccountRecord := func(ctx context.Context) (appuser.AccountRecord, error) {
//		if googleUser.Email == "" {
//			return nil, errors.New("Not implemented yet: Google did not provided appUser email")
//		}
//		userGoogle = models.NewUserAccountEntry(googleUser.ContactID)
//		data := userGoogle.DataStruct()
//		data.EmailData = appuser.NewEmailData(googleUser.Email)
//		data.ClientID = googleUser.ClientID
//		data.FederatedProvider = googleUser.FederatedProvider
//		data.FederatedIdentity = googleUser.FederatedIdentity
//		data.OwnedByUserWithID.AppUserID = appUserID
//		return &userGoogle, nil
//	}
//
//	if appUser, err = getOrCreateUserAccountRecordOnSignIn(
//		c,
//		"google",
//		appUserID,
//		getUserAccountRecordFromDB,
//		newUserAccountRecord,
//		clientInfo,
//	); err != nil {
//		return
//	}
//	return
//}

//func getOrCreateUserAccountRecordOnSignIn(
//	ctx context.Context,
//	provider string,
//	userID string,
//	getUserAccountRecordFromDB func(ctx context.Context) (appuser.AccountRecord, error),
//	newUserAccountRecord func(ctx context.Context) (appuser.AccountRecord, error),
//	clientInfo models.ClientInfo,
//) (
//	appUser models.AppUserOBSOLETE, err error,
//) {
//	logus.Debugf(ctx, "getOrCreateUserAccountRecordOnSignIn(provider=%v, userID=%d)", provider, userID)
//	var db dal.DB
//	if db, err = GetDatabase(ctx); err != nil {
//		return
//	}
//	var userAccount appuser.AccountRecord
//	err = db.RunReadwriteTransaction(ctx, func(ctx context.Context, tx dal.ReadwriteTransaction) (err error) {
//		if userAccount, err = getUserAccountRecordFromDB(ctx); err != nil {
//			if !dal.IsNotFound(err) {
//				// Technical error
//				return fmt.Errorf("failed to get user account record: %w", err)
//			}
//		}
//
//		userAccountRecord := dal.NewRecordWithData(dal.NewKeyWithID("UserEntry"+userAccount.Key().Provider, userAccount.Key().ContactID), userAccount.Data())
//
//		now := time.Now()
//
//		isNewUser := userID == ""
//
//		accountData := userAccount.Data()
//
//		updateUser := func() {
//			appUser.Data.SetLastLogin(now)
//			appUser.Data.SetLastLogin(now)
//			if !appUser.Data.EmailConfirmed && accountData.GetEmailConfirmed() {
//				appUser.Data.EmailConfirmed = true
//			}
//			names := accountData.GetNames()
//			if appUser.Data.FirstName == "" && names.FirstName != "" {
//				appUser.Data.FirstName = names.FirstName
//			}
//			if appUser.Data.LastName == "" && names.LastName != "" {
//				appUser.Data.LastName = names.LastName
//			}
//			if appUser.Data.Nickname == "" && names.NickName != "" {
//				appUser.Data.Nickname = names.NickName
//			}
//		}
//
//		if err == nil { // UserEntry account record found
//			uaRecordUserID := accountData.GetAppUserID()
//			if !isNewUser && uaRecordUserID != userID {
//				panic(fmt.Sprintf("Relinking of appUser accounts us not implemented yet => userAccount.GetAppUserIntID():%s != userID:%s", uaRecordUserID, userID))
//			}
//			if appUser, err = facade4userus.GetUserByID(c, tx, uaRecordUserID); err != nil {
//				if dal.IsNotFound(err) {
//					err = fmt.Errorf("record UserAccountEntry is referencing non existing appUser: %w", err)
//				}
//				return
//			}
//			accountData.SetLastLogin(now)
//			updateUser()
//
//			if err = tx.SetMulti(c, []dal.Record{userAccountRecord, appUser.Record}); err != nil {
//				return fmt.Errorf("failed to update UserEntry & UserFacebook with DtLastLogin: %w", err)
//			}
//			return
//		}
//
//		// UserAccountEntry record not found
//		// Lets create new UserAccountEntry record
//		if userAccount, err = newUserAccountRecord(c); err != nil {
//			return
//		}
//
//		if isNewUser {
//			//if i, ok := userAccount.(user.CreatedTimesSetter); ok {
//			//	i.SetCreatedTime(now)
//			//}
//		} else {
//			if appUser, err = facade4userus.GetUserByID(c, tx, userID); err != nil {
//				return
//			}
//		}
//
//		//if i, ok := userAccount.(user.UpdatedTimeSetter); ok {
//		//	i.SetUpdatedTime(now)
//		//}
//		accountData.SetLastLogin(now)
//
//		email := models.GetEmailID(accountData.GetEmailLowerCase())
//
//		if email == "" {
//			panic("Not implemented: userAccount.GetEmail() returned empty string")
//		}
//
//		var userEmail models.UserEmailEntry
//		if userEmail, err = dtdal.UserEmailEntry.GetUserEmailByID(c, tx, email); err != nil && !dal.IsNotFound(err) {
//			return // error
//		}
//
//		if dal.IsNotFound(err) { // UserEmailEntry record NOT found
//			userEmail := models.NewUserEmail(email, models.NewUserEmailData(0, true, provider))
//			userEmail.Data.CreatedAt = now
//
//			// We need to create new UserEntry entity
//			if isNewUser {
//				appUser = models.NewUserEntry(clientInfo)
//				appUser.Data.DtCreated = now
//			}
//			appUser.Data.AddAccount(userAccount.Key())       // No need to check for changed as new appUser
//			appUser.Data.AddAccount(userEmail.UserAccountEntry()) // No need to check for changed as new appUser
//			updateUser()
//
//			if isNewUser {
//				if appUser, err = dtdal.UserEntry.CreateUser(c, appUser.Data); err != nil {
//					return
//				}
//			} else if err = UserEntry.SaveUserOBSOLETE(c, tx, appUser); err != nil {
//				return
//			}
//
//			userAccount.(appuser.BelongsToUser).SetAppUserID(appUser.ContactID)
//			userEmail.Data.AppUserID = appUser.ContactID
//
//			if err = tx.SetMulti(c, []dal.Record{userAccountRecord, userEmail.Record}); err != nil {
//				return
//			}
//			return
//		} else { // UserEmailEntry record found
//			userAccount.(appuser.BelongsToUser).SetAppUserID(userEmail.Data.AppUserID) // No need to create a new appUser, link to existing
//			if !isNewUser && userEmail.Data.AppUserID != userID {
//				panic(fmt.Sprintf("Relinking of appUser accounts us not implemented yet => userEmail.AppUserID:%s != userID:%s", userEmail.Data.AppUserID, userID))
//			}
//
//			if isNewUser {
//				if appUser, err = facade4userus.GetUserByID(c, tx, userEmail.Data.AppUserID); err != nil {
//					if dal.IsNotFound(err) {
//						err = fmt.Errorf("record UserEmailEntry is referencing non existing UserEntry: %w", err)
//					}
//					return
//				}
//			}
//
//			if changed := userEmail.Data.AddProvider(provider); changed || !userEmail.Data.IsConfirmed {
//				userEmail.Data.IsConfirmed = true
//				if err = dtdal.UserEmailEntry.SaveUserEmail(c, tx, userEmail); err != nil {
//					return
//				}
//			}
//			appUser.Data.AddAccount(userAccount.Key())
//			updateUser()
//			if err = tx.SetMulti(c, []dal.Record{userAccountRecord, appUser.Record}); err != nil {
//				return fmt.Errorf("failed to create UserFacebook & update UserEntry: %w", err)
//			}
//			return
//		}
//	})
//	return
//}

//func (uf userFacade) GetOrCreateUserFacebookOnSignIn(
//	ctx context.Context,
//	appUserID int64,
//	fbAppOrPageID, fbUserOrPageScopeID, firstName, lastName string,
//	email string, isEmailConfirmed bool,
//	clientInfo models.ClientInfo,
//) (
//	userFacebook models.UserFacebook, appUser models.AppUserOBSOLETE, err error,
//) {
//	logus.Debugf(c, "GetOrCreateUserFacebookOnSignIn(firstName=%v, lastName=%v)", firstName, lastName)
//	if fbAppOrPageID == "" {
//		panic("fbAppOrPageID is empty string")
//	}
//	if fbAppOrPageID == "" {
//		panic("fbUserOrPageScopeID is empty string")
//	}
//
//	updateNames := func(entity *models.UserFacebookData) {
//		if firstName != "" && userFacebook.Data.FirstName != firstName {
//			userFacebook.Data.FirstName = firstName
//		}
//		if lastName != "" && userFacebook.Data.LastName != lastName {
//			userFacebook.Data.LastName = lastName
//		}
//	}
//
//	getUserAccountRecordFromDB := func(ctx context.Context) (user.AccountRecord, error) {
//		if userFacebook, err = dtdal.UserFacebook.GetFbUserByFbID(ctx, fbAppOrPageID, fbUserOrPageScopeID); err != nil {
//			return &userFacebook, err
//		}
//		updateNames(userFacebook.Data)
//		return &userFacebook, err
//	}
//
//	newUserAccountRecord := func(ctx context.Context) (user.AccountRecord, error) {
//		userFacebook = models.UserFacebook{
//			FbAppOrPageID:       fbAppOrPageID,
//			FbUserOrPageScopeID: fbUserOrPageScopeID,
//			Data: &models.UserFacebookData{
//				Email: email,
//				Names: user.Names{
//					FirstName: firstName,
//					LastName:  lastName,
//				},
//				EmailIsConfirmed: isEmailConfirmed,
//				OwnedByUserWithID: user.OwnedByUserWithID{
//					AppUserIntID: appUserID,
//					AppUserID:    strconv.FormatInt(appUserID, 10),
//				},
//			},
//		}
//		updateNames(userFacebook.Data)
//		return &userFacebook, nil
//	}
//	if appUser, err = getOrCreateUserAccountRecordOnSignIn(
//		c,
//		"fb",
//		appUserID,
//		getUserAccountRecordFromDB,
//		newUserAccountRecord,
//		clientInfo,
//	); err != nil {
//		return
//	}
//	return
//}
