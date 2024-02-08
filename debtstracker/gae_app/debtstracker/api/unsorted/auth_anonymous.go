package unsorted

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"net/http"
)

func HandleSignUpAnonymously(c context.Context, w http.ResponseWriter, r *http.Request) {
	if user, err := dtdal.User.CreateAnonymousUser(c); err != nil {
		api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
	} else {
		if err = SaveUserAgent(c, user.ID, r.UserAgent()); err != nil {
			api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			return
		}
		ReturnToken(c, w, user.ID, true, false)
	}
}

func HandleSignInAnonymous(c context.Context, w http.ResponseWriter, r *http.Request) {
	userID := r.PostFormValue("user")
	if userID == "" {
		api.ErrorAsJson(c, w, http.StatusBadRequest, errors.New("required parameter user is empty"))
		return
	}

	userEntity, err := facade.User.GetUserByID(c, nil, userID)

	if err != nil {
		if dal.IsNotFound(err) {
			api.ErrorAsJson(c, w, http.StatusBadRequest, err)
		} else {
			api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
		}
		return
	}

	if userEntity.Data.IsAnonymous {
		if err = SaveUserAgent(c, userID, r.UserAgent()); err != nil {
			api.ErrorAsJson(c, w, http.StatusInternalServerError, err)
			return
		}
		ReturnToken(c, w, userID, false, false)
	} else {
		api.ErrorAsJson(c, w, http.StatusForbidden, errors.New("User is not anonymous."))
	}
}

//func handleLinkOneSignal(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
//	_, err := dtdal.UserOneSignal.SaveUserOneSignal(c, authInfo.UserID, r.PostFormValue("OneSignalUserID"))
//	if err != nil {
//		ErrorAsJson(c, w, http.StatusInternalServerError, err)
//	}
//}
