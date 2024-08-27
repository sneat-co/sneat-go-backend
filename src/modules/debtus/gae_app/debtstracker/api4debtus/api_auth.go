package api4debtus

import (
	"context"
	"errors"
	"github.com/sneat-co/sneat-go-backend/src/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp"
	"github.com/strongo/validation"
	"net/http"
)

type AuthHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo)

type AuthHandlerWithUser func(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo, user dbo4userus.UserEntry)

func AuthOnly(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if authInfo, _, err := token4auth.Authenticate(w, r, true); err == nil {
			handler(ctx, w, r, authInfo)
		} else {
			logus.Warningf(ctx, "Failed to authenticate: %v", err)
		}
	}
}

func AuthOnlyWithUser(handler AuthHandlerWithUser) strongoapp.HttpHandlerWithContext {
	return AuthOnly(func(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
		var userID string

		if userID = GetUserID(ctx, w, r, authInfo); userID == "" {
			logus.Warningf(ctx, "userID is empty")
			return
		}

		//user, err := dal4userus.GetUserByID(ctx, nil, userID)
		var user dbo4userus.UserEntry
		var err error

		if HasError(ctx, w, err, models4debtus.AppUserKind, userID, http.StatusInternalServerError) {
			return
		}
		handler(ctx, w, r, authInfo, user)
	})
}

func OptionalAuth(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		authInfo, _, _ := token4auth.Authenticate(w, r, false)
		if authInfo.UserID == "" {
			logus.Debugf(ctx, "OptionalAuth(), anonymous")
		} else {
			logus.Debugf(ctx, "OptionalAuth(), userID=%s", authInfo.UserID)
		}
		handler(ctx, w, r, authInfo)
	}
}

func AdminOnly(handler AuthHandler) strongoapp.HttpHandlerWithContext {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if authInfo, _, err := token4auth.Authenticate(w, r, true); err == nil {
			if !authInfo.IsAdmin {
				logus.Debugf(ctx, "Not admin!")
				//hashedWriter.WriteHeader(http.StatusForbidden)
				//return
			}
			handler(ctx, w, r, authInfo)
		} else {
			logus.Errorf(ctx, "Failed to authenticate: %v", err)
		}
	}
}

func IsAdmin(email string) bool {
	return email == "alexander.trakhimenok@gmail.com"
}

func ReturnToken(ctx context.Context, w http.ResponseWriter, userID, issuer string, isNewUser, isAdmin bool) {
	if isAdmin {
		apicore.ReturnError(ctx, w, nil, validation.NewBadRequestError(errors.New("issuing admin token is not implemented yet")))
		return
	}
	token, err := token4auth.IssueFirebaseAuthToken(ctx, userID, issuer)
	if err != nil {
		apicore.ReturnError(ctx, w, nil, err)
		return
	}
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
