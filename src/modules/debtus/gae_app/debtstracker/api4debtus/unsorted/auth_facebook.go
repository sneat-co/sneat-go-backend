package unsorted

import (
	"github.com/sneat-co/sneat-go-backend/src/auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/strongo/logus"
	"net/http"

	"context"
)

func HandleSignedWithFacebook(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	logus.Debugf(c, "api4debtus.HandleSignedWithFacebook()")
	fbUserID := r.PostFormValue("fbUserID")
	fbAppID := r.PostFormValue("fbAppID")
	if fbUserID == "" {
		api4debtus.BadRequestMessage(c, w, "fbUserID is missed")
		return
	}
	if fbAppID == "" {
		api4debtus.BadRequestMessage(c, w, "fbAppID is missed")
		return
	}
	user, isNewUser, _, _, _, err := signInFbUser(c, fbAppID, fbUserID, r, authInfo)
	if err != nil {
		authWriteResponseForAuthFailed(c, w, err)
		return
	}
	authWriteResponseForUser(c, w, user, isNewUser)
}
