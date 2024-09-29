package api4auth

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/common4all"
	"github.com/strongo/logus"
	"net/http"

	"context"
)

func HandleSignedWithFacebook(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	logus.Debugf(ctx, "api4debtus.HandleSignedWithFacebook()")
	fbUserID := r.PostFormValue("fbUserID")
	fbAppID := r.PostFormValue("fbAppID")
	if fbUserID == "" {
		common4all.BadRequestMessage(ctx, w, "fbUserID is missed")
		return
	}
	if fbAppID == "" {
		common4all.BadRequestMessage(ctx, w, "fbAppID is missed")
		return
	}
	user, isNewUser, _, _, _, err := signInFbUser(ctx, fbAppID, fbUserID, r, authInfo)
	if err != nil {
		authWriteResponseForAuthFailed(ctx, w, err)
		return
	}
	authWriteResponseForUser(ctx, w, r, user, "facebook", isNewUser)
}
