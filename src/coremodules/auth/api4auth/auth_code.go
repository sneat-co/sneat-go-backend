package api4auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/facade4anybot"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/models4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/strongo/logus"
	"net/http"
	"strconv"
)

// TODO: Obsolete - migrate to HandleSignInWithPin
func HandleSignInWithCode(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	code := r.PostFormValue("code")
	if code == "" {
		api4debtus.BadRequestMessage(ctx, w, "Missing required attribute: code")
		return
	}
	if loginCode, err := strconv.Atoi(code); err != nil {
		api4debtus.BadRequestMessage(ctx, w, "Parameter code is not an integer")
		return
	} else if loginCode == 0 {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Login code should not be 0."))
		return
	} else {
		if userID, err := dtdal.LoginCode.ClaimLoginCode(ctx, loginCode); err != nil {
			switch err {
			case models4auth.ErrLoginCodeExpired:
				_, _ = w.Write([]byte("expired"))
			case models4auth.ErrLoginCodeAlreadyClaimed:
				_, _ = w.Write([]byte("claimed"))
			default:
				err = fmt.Errorf("failed to claim code: %w", err)
				api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			}
		} else {
			if authInfo.UserID != "" && userID != authInfo.UserID {
				logus.Warningf(ctx, "userID:%s != authInfo.AppUserIntID:%s", userID, authInfo.UserID)
			}
			ReturnToken(ctx, w, r, userID, r.Referer())
			return
		}
	}
}

func HandleSignInWithPin(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	loginID, err := common4debtus.DecodeIntID(r.PostFormValue("loginID"))
	if err != nil {
		api4debtus.BadRequestError(ctx, w, fmt.Errorf("parameter 'loginID' is not an integer: %w", err))
		return
	}

	if loginCode, err := strconv.ParseInt(r.PostFormValue("loginPin"), 10, 32); err != nil {
		api4debtus.BadRequestMessage(ctx, w, "Parameter 'loginCode' is not an integer")
		return
	} else if loginCode == 0 {
		api4debtus.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("Parameter 'loginCode' should not be 0."))
		return
	} else {
		if userID, err := facade4anybot.AuthFacade.SignInWithPin(ctx, loginID, int32(loginCode)); err != nil {
			switch err {
			case facade4anybot.ErrLoginExpired:
				_, _ = w.Write([]byte("expired"))
			case facade4anybot.ErrLoginAlreadySigned:
				_, _ = w.Write([]byte("claimed"))
			default:
				err = fmt.Errorf("failed to claim loginCode: %w", err)
				api4debtus.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			}
		} else {
			if authInfo.UserID != "" && userID != authInfo.UserID {
				logus.Warningf(ctx, "userID:%s != authInfo.AppUserIntID:%s", userID, authInfo.UserID)
			}
			ReturnToken(ctx, w, r, userID, r.Referer())
		}
	}
}
