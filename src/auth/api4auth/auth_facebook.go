package api4auth

import (
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/strongo/logus"
	"net/http"

	"context"
)

func HandleSignedWithFacebook(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	logus.Debugf(ctx, "api4debtus.HandleSignedWithFacebook()")
	fbUserID := r.PostFormValue("fbUserID")
	fbAppID := r.PostFormValue("fbAppID")
	if fbUserID == "" {
		api4debtus.BadRequestMessage(ctx, w, "fbUserID is missed")
		return
	}
	if fbAppID == "" {
		api4debtus.BadRequestMessage(ctx, w, "fbAppID is missed")
		return
	}
	user, isNewUser, _, _, _, err := signInFbUser(ctx, fbAppID, fbUserID, r, authInfo)
	if err != nil {
		authWriteResponseForAuthFailed(ctx, w, err)
		return
	}
	authWriteResponseForUser(ctx, w, user, "facebook", isNewUser)
}
