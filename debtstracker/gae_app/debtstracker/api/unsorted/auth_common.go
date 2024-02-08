package unsorted

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
	"github.com/strongo/strongoapp/appuser"
	"net/http"
)

func HandleDisconnect(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	provider := r.URL.Query().Get("provider")

	var err error
	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		return
	}
	if err := db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		appUser, err := facade.User.GetUserByID(c, tx, authInfo.UserID)
		if err != nil {
			return err
		}

		changed := false

		deleteFbUser := func(userAccount appuser.AccountKey) error {
			if userFb, err := dtdal.UserFacebook.GetFbUserByFbID(c, userAccount.App, userAccount.ID); err != nil {
				if !dal.IsNotFound(err) {
					return err
				}
			} else if fbUserAppID := userFb.FbUserData().GetAppUserID(); fbUserAppID == appUser.ID {
				if err = dtdal.UserFacebook.DeleteFbUser(c, userAccount.App, userAccount.ID); err != nil {
					return err
				}
			} else {
				log.Warningf(c, "TODO: Handle case if userFb.AppUserIntID:%s != appUser.ID:%d", fbUserAppID, appUser.ID)
			}
			return nil
		}

		if !models.IsKnownUserAccountProvider(provider) {
			api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Unknown provider: "+provider))
			return nil
		}
		if !appUser.Data.HasAccount(provider, "") {
			return nil
		}
		var userAccount *appuser.AccountKey
		switch provider {
		case "google":
			if userAccount, err = appUser.Data.GetAccount("google", ""); err != nil {
				return err
			} else if userAccount != nil {
				if userGoogle, err := dtdal.UserGoogle.GetUserGoogleByID(c, userAccount.ID); err != nil {
					if !dal.IsNotFound(err) {
						return err
					}
				} else if userGoogle.Data().GetAppUserID() == appUser.ID {
					userGoogle.Data().SetAppUserID("")
					if err = dtdal.UserGoogle.DeleteUserGoogle(c, userGoogle.ID); err != nil {
						return err
					}
				} else {
					log.Warningf(c, "TODO: Handle case if userGoogle.AppUserIntID:%d != appUser.ID:%d", userGoogle.Data().GetAppUserID(), appUser.ID)
				}
				_ = appUser.Data.RemoveAccount(*userAccount)
				changed = true
			}
		case "fb":
			if userAccount, err = appUser.Data.GetAccount("facebook", ""); err != nil {
				return err
			} else if userAccount != nil {
				if err = deleteFbUser(*userAccount); err != nil {
					return err
				}
				_ = appUser.Data.RemoveAccount(*userAccount)
				changed = true
			}
		case "fbm":
			if userAccount, err = appUser.Data.GetAccount("facebook", ""); err != nil {
				return err
			} else if userAccount != nil {
				if err = deleteFbUser(*userAccount); err != nil {
					return err
				}
				_ = appUser.Data.RemoveAccount(*userAccount)
				changed = true
			}
		default:
		}

		if changed {
			if err = facade.User.SaveUser(c, tx, appUser); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
	}
}
