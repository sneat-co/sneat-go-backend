package gaedal

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

func NewUserGoogleKey(id string) *dal.Key {
	return dal.NewKeyWithID(models.UserGoogleCollection, id)
}

type UserGoogleDalGae struct {
}

func NewUserGoogleDalGae() UserGoogleDalGae {
	return UserGoogleDalGae{}
}

func (UserGoogleDalGae) GetUserGoogleByID(c context.Context, googleUserID string) (userGoogle models.UserAccount, err error) {
	//userGoogle.ID = googleUserID
	//userGoogle.Data = new(models.UserGoogleData)
	//if err = gaedb.Get(c, NewUserGoogleKey(googleUserID), userGoogle.Data); err != nil {
	//	if err == datastore.ErrNoSuchEntity {
	//		err = dal.ErrRecordNotFound
	//	}
	//	return
	//}
	err = errors.New("not implemented")
	return
}

func (UserGoogleDalGae) DeleteUserGoogle(c context.Context, googleUserID string) (err error) {
	//if err = gaedb.Delete(c, NewUserGoogleKey(googleUserID)); err != nil {
	//	return
	//}
	return errors.New("not implemented")
}

func (UserGoogleDalGae) SaveUserGoogle(c context.Context, userGoogle models.UserAccount) (err error) {
	//if _, err = gaedb.Put(c, NewUserGoogleKey(userGoogle.ID), userGoogle.Data); err != nil {
	//	return
	//}
	return errors.New("not implemented")
}

// TODO: Obsolete!
//func (UserGoogleDalGae) CreateUserGoogle(c context.Context, user user.User, appUserID int64, onSignIn bool, userAgent, remoteAddr string) (entity *models.UserGoogleData, isNewGoogleUser, isNewAppUser bool, err error) {
//	err = dtdal.DB.RunInTransaction(c, func(tc context.Context) (err error) {
//		key := NewUserGoogleKey(tc, user.ID)
//		entity = new(models.UserGoogleData)
//
//		if err = gaedb.Get(tc, key, entity); err == nil {
//			if onSignIn {
//				entity.LastSignIn = time.Now()
//
//				if appUserID != 0 && entity.AppUserIntID != appUserID { // Reconnect Google account to different user
//
//					if entity.AppUserIntID == 0 {
//						if appUser, err := facade.User.GetUserByID(c, appUserID); err != nil {
//							return err
//						} else /* if appUser.GoogleUniqueUserID == "" */ {
//							appUser.GoogleUniqueUserID = user.ID
//							if err = facade.User.SaveUser(c, appUser); err != nil {
//								return err
//							}
//						} // TODO: Handle case when appUser.GoogleUniqueUserID is not empty
//					} else {
//						oldUser := models.AppUser{ID: entity.AppUserIntID}
//						newUser := models.AppUser{ID: appUserID}
//
//						if err = dtdal.DB.GetMulti(c, []db.EntityHolder{&oldUser, &newUser}); err != nil {
//							return
//						}
//
//						oldUser.GoogleUniqueUserID = ""
//						newUser.GoogleUniqueUserID = user.ID
//
//						if err = dtdal.DB.UpdateMulti(c, []db.EntityHolder{&oldUser, &newUser}); err != nil {
//							return
//						}
//					}
//					entity.AppUserIntID = appUserID
//				}
//
//				if _, err = gaedb.Put(tc, key, entity); err != nil {
//					err = errors.Wrap(err, "Failed to save google user")
//					return
//				}
//			}
//			return
//		} else  if err != datastore.ErrNoSuchEntity {
//			err = errors.Wrapf(err, "Failed to get google user entity by key=%v", key)
//			return
//		}
//
//		isNewGoogleUser = true
//		now := time.Now()
//		entity = &models.UserGoogleData{
//			LastSignIn: now,
//			User:       user,
//			OwnedByUserWithIntID: user.OwnedByUserWithIntID{
//				AppUserIntID: appUserID,
//				DtCreated: now,
//			},
//		}
//
//		if entity.AppUserIntID != 0 {
//			if user, err := facade.User.GetUserByID(c, entity.AppUserIntID); err != nil {
//				return err
//			} else if user.GoogleUniqueUserID != entity.ID {
//				if user.GoogleUniqueUserID != "" {
//					log.Warningf(c, "TODO: Handle case when connect with to user with different linked Google ID")
//				}
//				user.GoogleUniqueUserID = entity.ID
//				if err = facade.User.SaveUser(c, user); err != nil {
//					return err
//				}
//			}
//		} else {
//			emailLowCase := strings.ToLower(user.Email)
//			query := datastore.NewQuery(models.AppUserKind).Filter("EmailAddress = ", emailLowCase).Limit(2)
//			var (
//				appUserKeys []*datastore.Key
//				appUsers    []models.DebutsAppUserDataOBSOLETE
//			)
//			if appUserKeys, err = query.GetAll(c, &appUsers); err != nil {
//				err = errors.Wrap(err, "Failed to load users by email")
//				return
//			}
//			switch len(appUserKeys) {
//			case 1:
//				entity.AppUserIntID = appUserKeys[0].IntegerID()
//			case 0:
//				query = datastore.NewQuery(models.UserGoogleCollection).Filter("Email =", user.Email).Limit(2)
//				var (
//					googleUserKeys []*datastore.Key
//					googleUsers    []models.UserGoogleData
//				)
//				if googleUserKeys, err = query.GetAll(c, &googleUsers); err != nil {
//					err = errors.Wrap(err, "Failed to load google users by email")
//					return
//				}
//				switch len(googleUserKeys) {
//				case 1:
//					panic("TODO: We need to handle situation when user changed email and that email was linked to another google account")
//				case 2:
//					err = fmt.Errorf("Found > 1 google users for email=%v, %v", user.Email, googleUserKeys)
//					return
//				}
//
//				isNewAppUser = true
//				appUserKey := datastore.NewIncompleteKey(tc, models.AppUserKind, nil)
//				if strings.Index(remoteAddr, ":") >= 0 {
//					remoteAddr = strings.Split(remoteAddr, ":")[0]
//				}
//				appUser := models.DebutsAppUserDataOBSOLETE{
//					GoogleUniqueUserID: user.ID,
//					DtCreated:          now,
//					LastUserAgent:      userAgent,
//					LastUserIpAddress:  remoteAddr,
//					ContactDetails: models.ContactDetails{
//						EmailContact: models.EmailContact{
//							EmailAddress:         emailLowCase,
//							EmailAddressOriginal: user.Email,
//							EmailConfirmed:       true,
//						},
//					},
//				}
//				if appUserKey, err = gaedb.Put(tc, appUserKey, &appUser); err != nil {
//					err = errors.Wrap(err, "Failed to save app use entity")
//					return
//				}
//				entity.AppUserIntID = appUserKey.IntegerID()
//			default: // len(appUserKeys) > 1
//				err = fmt.Errorf("Found > 1 users for email=%v, %v", emailLowCase, appUserKeys)
//				return
//			}
//		}
//
//		if _, err = gaedb.Put(tc, key, entity); err != nil {
//			err = errors.Wrap(err, "Failed to save google use entity")
//			return
//		}
//		return
//	}, dtdal.CrossGroupTransaction)
//	return
//}
