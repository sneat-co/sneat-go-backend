package api4auth

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/const4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp/appuser"
	"net/http"
)

func HandleDisconnect(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	provider := r.URL.Query().Get("provider")

	userCtx := facade.NewUserContext(authInfo.UserID)

	if err := dal4userus.RunUserWorker(ctx, userCtx, true, func(ctx context.Context, tx dal.ReadwriteTransaction, params *dal4userus.UserWorkerParams) error {
		appUser, err := dal4userus.GetUserByID(ctx, tx, authInfo.UserID)
		if err != nil {
			return err
		}

		changed := false

		deleteFbUser := func(userAccount appuser.AccountKey) error {
			if userFb, err := facade4auth.UserFacebook.GetFbUserByFbID(ctx, userAccount.App, userAccount.ID); err != nil {
				if !dal.IsNotFound(err) {
					return err
				}
			} else if fbUserAppID := userFb.FbUserData().GetAppUserID(); fbUserAppID == appUser.ID {
				if err = facade4auth.UserFacebook.DeleteFbUser(ctx, userAccount.App, userAccount.ID); err != nil {
					return err
				}
			} else {
				logus.Warningf(ctx, "TODO: Handle case if userFb.AppUserIntID:%s != appUser.ContactID:%s", fbUserAppID, appUser.ID)
			}
			return nil
		}

		if !const4userus.IsKnownUserAccountProvider(provider) {
			api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Unknown provider: "+provider))
			return nil
		}
		if !appUser.Data.HasAccount(provider, "") {
			return nil
		}
		var userAccount *appuser.AccountKey
		switch provider {
		case const4userus.GoogleAuthProvider:
			if userAccount, err = appUser.Data.GetAccount(provider, ""); err != nil {
				return err
			} else if userAccount != nil {
				if userGoogle, err := dtdal.UserGoogle.GetUserGoogleByID(ctx, userAccount.ID); err != nil {
					if !dal.IsNotFound(err) {
						return err
					}
				} else if userGoogle.Data.GetAppUserID() == appUser.ID {
					userGoogle.Data.SetAppUserID("")
					if err = dtdal.UserGoogle.DeleteUserGoogle(ctx, userGoogle.ID); err != nil {
						return err
					}
				} else {
					logus.Warningf(ctx, "TODO: Handle case if userGoogle.AppUserIntID:%s != appUser.ContactID:%s", userGoogle.Data.GetAppUserID(), appUser.ID)
				}
				_ = appUser.Data.RemoveAccount(*userAccount)
				changed = true
			}
		case const4userus.FacebookAuthProvider:
			if userAccount, err = appUser.Data.GetAccount("facebook", ""); err != nil {
				return err
			} else if userAccount != nil {
				if err = deleteFbUser(*userAccount); err != nil {
					return err
				}
				_ = appUser.Data.RemoveAccount(*userAccount)
				changed = true
			}
		case const4userus.FacebookMessengerAuthProvider:
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
			appUser.Record.MarkAsChanged()
		}
		return nil
	}); err != nil {
		api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	}
}
