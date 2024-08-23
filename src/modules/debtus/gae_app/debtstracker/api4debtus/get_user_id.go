package api4debtus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"net/http"
)

func GetUserID(_ context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) (userID string) {
	userID = authInfo.UserID

	if stringID := r.URL.Query().Get("user"); stringID != "" {
		if !authInfo.IsAdmin && userID != authInfo.UserID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
	return
}
