package api4debtus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"net/http"
)

type AuthHandler func(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo)

type AuthHandlerWithUser func(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo, user dbo4userus.UserEntry)

func AuthOnly(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		if authInfo, _, err := auth.Authenticate(w, r, true); err == nil {
			handler(c, w, r, authInfo)
		} else {
			logus.Warningf(c, "Failed to authenticate: %v", err)
		}
	}
}

func AuthOnlyWithUser(handler AuthHandlerWithUser) strongoapp.HttpHandlerWithContext {
	return AuthOnly(func(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
		var userID string

		if userID = GetUserID(c, w, r, authInfo); userID == "" {
			logus.Warningf(c, "userID is empty")
			return
		}

		//user, err := dal4userus.GetUserByID(c, nil, userID)
		var user dbo4userus.UserEntry
		var err error

		if HasError(c, w, err, models4debtus.AppUserKind, userID, http.StatusInternalServerError) {
			return
		}
		handler(c, w, r, authInfo, user)
	})
}

func OptionalAuth(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		authInfo, _, _ := auth.Authenticate(w, r, false)
		if authInfo.UserID == "" {
			logus.Debugf(c, "OptionalAuth(), anonymous")
		} else {
			logus.Debugf(c, "OptionalAuth(), userID=%s", authInfo.UserID)
		}
		handler(c, w, r, authInfo)
	}
}

func AdminOnly(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) {
		if authInfo, _, err := auth.Authenticate(w, r, true); err == nil {
			if !authInfo.IsAdmin {
				logus.Debugf(c, "Not admin!")
				//hashedWriter.WriteHeader(http.StatusForbidden)
				//return
			}
			handler(c, w, r, authInfo)
		} else {
			logus.Errorf(c, "Failed to authenticate: %v", err)
		}
	}
}

func IsAdmin(email string) bool {
	return email == "alexander.trakhimenok@gmail.com"
}

func ReturnToken(_ context.Context, w http.ResponseWriter, userID string, isNewUser, isAdmin bool) {
	token := auth.IssueToken(userID, "api4debtus", isAdmin)
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Content-Type", "application/json")
	_, _ = w.Write([]byte("{"))
	if isNewUser {
		_, _ = w.Write([]byte(`"isNewUser":true,`))
	}
	_, _ = w.Write([]byte(`"token":"`))
	_, _ = w.Write([]byte(token))
	_, _ = w.Write([]byte(`"}`))
}
