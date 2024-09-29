package api4auth

import (
	"context"
	"errors"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/unsorted4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/common4all"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/userus/facade4userus"
	"net/http"
)

func HandleSignUpAnonymously(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if user, err := unsorted4auth.User.CreateAnonymousUser(ctx); err != nil {
		common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
	} else {
		if _, err = facade4userus.SaveUserBrowser(ctx, user.ID, r.UserAgent()); err != nil {
			common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			return
		}
		ReturnToken(ctx, w, r, user.ID, "")
	}
}

func HandleSignInAnonymous(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	userID := r.PostFormValue("user")
	if userID == "" {
		common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, errors.New("required parameter user is empty"))
		return
	}

	userEntity, err := dal4userus.GetUserByID(ctx, nil, userID)

	if err != nil {
		if dal.IsNotFound(err) {
			common4all.ErrorAsJson(ctx, w, http.StatusBadRequest, err)
		} else {
			common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
		}
		return
	}

	if userEntity.Data.IsAnonymous {
		if _, err = facade4userus.SaveUserBrowser(ctx, userID, r.UserAgent()); err != nil {
			common4all.ErrorAsJson(ctx, w, http.StatusInternalServerError, err)
			return
		}
		ReturnToken(ctx, w, r, userID, r.Referer())
	} else {
		common4all.ErrorAsJson(ctx, w, http.StatusForbidden, errors.New("User is not anonymous."))
	}
}

//func handleLinkOneSignal(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
//	_, err := facade4userus.SaveUserOneSignal(ctx, authInfo.UserID, r.PostFormValue("OneSignalUserID"))
//	if err != nil {
//		ErrorAsJson(c, w, http.StatusInternalServerError, err)
//	}
//}
