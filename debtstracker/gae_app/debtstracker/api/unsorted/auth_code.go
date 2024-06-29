package unsorted

import (
	"fmt"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/strongo/logus"
	"net/http"
	"strconv"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

// TODO: Obsolete - migrate to HandleSignInWithPin
func HandleSignInWithCode(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	code := r.PostFormValue("code")
	if code == "" {
		api.BadRequestMessage(c, w, "Missing required attribute: code")
		return
	}
	if loginCode, err := strconv.Atoi(code); err != nil {
		api.BadRequestMessage(c, w, "Parameter code is not an integer")
		return
	} else if loginCode == 0 {
		api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Login code should not be 0."))
		return
	} else {
		if userID, err := dtdal.LoginCode.ClaimLoginCode(c, loginCode); err != nil {
			switch err {
			case models.ErrLoginCodeExpired:
				_, _ = w.Write([]byte("expired"))
			case models.ErrLoginCodeAlreadyClaimed:
				_, _ = w.Write([]byte("claimed"))
			default:
				err = fmt.Errorf("failed to claim code: %w", err)
				api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			}
		} else {
			if authInfo.UserID != "" && userID != authInfo.UserID {
				logus.Warningf(c, "userID:%s != authInfo.AppUserIntID:%s", userID, authInfo.UserID)
			}
			ReturnToken(c, w, userID, false, false)
			return
		}
	}
}

func HandleSignInWithPin(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	loginID, err := common.DecodeIntID(r.PostFormValue("loginID"))
	if err != nil {
		api.BadRequestError(c, w, fmt.Errorf("parameter 'loginID' is not an integer: %w", err))
		return
	}

	if loginCode, err := strconv.ParseInt(r.PostFormValue("loginPin"), 10, 32); err != nil {
		api.BadRequestMessage(c, w, "Parameter 'loginCode' is not an integer")
		return
	} else if loginCode == 0 {
		api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Parameter 'loginCode' should not be 0."))
		return
	} else {
		if userID, err := facade.AuthFacade.SignInWithPin(c, loginID, int32(loginCode)); err != nil {
			switch err {
			case facade.ErrLoginExpired:
				_, _ = w.Write([]byte("expired"))
			case facade.ErrLoginAlreadySigned:
				_, _ = w.Write([]byte("claimed"))
			default:
				err = fmt.Errorf("failed to claim loginCode: %w", err)
				api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			}
		} else {
			if authInfo.UserID != "" && userID != authInfo.UserID {
				logus.Warningf(c, "userID:%s != authInfo.AppUserIntID:%s", userID, authInfo.UserID)
			}
			ReturnToken(c, w, userID, false, false)
		}
	}
}
