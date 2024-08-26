package token4auth

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/sneat-co/sneat-go-backend/src/core/facade2firebase"
)

func IssueFirebaseAuthToken(ctx context.Context, userID string, issuer string) (token string, err error) {

	var fbAuthClient *auth.Client
	if fbAuthClient, err = facade2firebase.GetFirebaseAuthClient(ctx); err != nil {
		return
	}
	claims := map[string]interface{}{}
	if issuer != "" {
		claims["issuer"] = issuer
	}
	return fbAuthClient.CustomTokenWithClaims(ctx, userID, claims)
}
