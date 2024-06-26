package unsorted

import (
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/api"
	"github.com/strongo/logus"
	"net/http"

	"context"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
)

func HandleSignedWithFacebook(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	logus.Debugf(c, "api.HandleSignedWithFacebook()")
	fbUserID := r.PostFormValue("fbUserID")
	fbAppID := r.PostFormValue("fbAppID")
	if fbUserID == "" {
		api.BadRequestMessage(c, w, "fbUserID is missed")
		return
	}
	if fbAppID == "" {
		api.BadRequestMessage(c, w, "fbAppID is missed")
		return
	}
	user, isNewUser, _, _, _, err := signInFbUser(c, fbAppID, fbUserID, r, authInfo)
	if err != nil {
		authWriteResponseForAuthFailed(c, w, err)
		return
	}
	authWriteResponseForUser(c, w, user, isNewUser)
}
