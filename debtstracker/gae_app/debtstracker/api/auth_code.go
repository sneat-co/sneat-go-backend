package api

import (
	"fmt"
	"net/http"
	"strconv"

	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

// TODO: Obsolete - migrate to handleSignInWithPin
func handleSignInWithCode(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	code := r.PostFormValue("code")
	if code == "" {
		BadRequestMessage(c, w, "Missing required attribute: code")
		return
	}
	if loginCode, err := strconv.Atoi(code); err != nil {
		BadRequestMessage(c, w, "Parameter code is not an integer")
		return
	} else if loginCode == 0 {
		ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Login code should not be 0."))
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
				ErrorAsJson(c, w, http.StatusInternalServerError, err)
			}
		} else {
			if authInfo.UserID != "" && userID != authInfo.UserID {
				log.Warningf(c, "userID:%v != authInfo.AppUserIntID:%v", userID, authInfo.UserID)
			}
			ReturnToken(c, w, userID, false, false)
			return
		}
	}
}

func handleSignInWithPin(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	loginID, err := common.DecodeIntID(r.PostFormValue("loginID"))
	if err != nil {
		BadRequestError(c, w, fmt.Errorf("parameter 'loginID' is not an integer: %w", err))
		return
	}

	if loginCode, err := strconv.ParseInt(r.PostFormValue("loginPin"), 10, 32); err != nil {
		BadRequestMessage(c, w, "Parameter 'loginCode' is not an integer")
		return
	} else if loginCode == 0 {
		ErrorAsJson(c, w, http.StatusBadRequest, errors.New("Parameter 'loginCode' should not be 0."))
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
				ErrorAsJson(c, w, http.StatusInternalServerError, err)
			}
		} else {
			if authInfo.UserID != "" && userID != authInfo.UserID {
				log.Warningf(c, "userID:%v != authInfo.AppUserIntID:%v", userID, authInfo.UserID)
			}
			ReturnToken(c, w, userID, false, false)
		}
	}
}
