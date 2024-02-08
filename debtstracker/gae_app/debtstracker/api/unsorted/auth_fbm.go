package unsorted

import (
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"net/http"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/strongo/log"
)

func HandleSignInWithFbm(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	log.Debugf(c, "api.HandleSignInWithFbm()")

	threadID := r.PostFormValue("tid")
	threadType := r.PostFormValue("thread_type")
	pageScopedID := r.PostFormValue("psid")
	fbAppID := r.PostFormValue("fbAppID")

	if fbAppID == "" {
		api.BadRequestMessage(c, w, "Missing fbAppID")
		return
	}
	if threadID == "" {
		api.BadRequestMessage(c, w, "Missing tid")
		return
	}
	if threadType == "" {
		api.BadRequestMessage(c, w, "Missing thread_type")
		return
	}
	if pageScopedID == "" {
		api.BadRequestMessage(c, w, "Missing psid")
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
