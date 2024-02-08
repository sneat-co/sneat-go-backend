package api

import (
	"net/http"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/strongo/log"
)

func handleSignInWithFbm(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	log.Debugf(c, "api.handleSignInWithFbm()")

	threadID := r.PostFormValue("tid")
	threadType := r.PostFormValue("thread_type")
	pageScopedID := r.PostFormValue("psid")
	fbAppID := r.PostFormValue("fbAppID")

	if fbAppID == "" {
		BadRequestMessage(c, w, "Missing fbAppID")
		return
	}
	if threadID == "" {
		BadRequestMessage(c, w, "Missing tid")
		return
	}
	if threadType == "" {
		BadRequestMessage(c, w, "Missing thread_type")
		return
	}
	if pageScopedID == "" {
		BadRequestMessage(c, w, "Missing psid")
		return
	}

	// TODO: Log FbApp Code & FbPage Code (e.g. fbAppID=12345 => code=DebtsTracker)
	log.Debugf(c, "FbmContext: thread_type=%v, tid=%v, psid=%v", threadType, threadID, pageScopedID)

	user, isNewUser, _, _, _, err := signInFbUser(c, fbAppID, pageScopedID, r, authInfo)
	if err != nil {
		authWriteResponseForAuthFailed(c, w, err)
		return
	}

	authWriteResponseForUser(c, w, user, isNewUser)
}
