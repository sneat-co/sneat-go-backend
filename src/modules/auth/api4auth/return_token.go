package api4auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/token4auth"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/strongo/validation"
	"net/http"
)

type TokenClaim func(claims *TokenClaims)

type TokenClaims struct {
	isAdmin bool
}

func (t *TokenClaims) IsAdmin() bool {
	return t.isAdmin
}

func IsAdminClaim() func(claim *TokenClaims) {
	return func(claim *TokenClaims) {
		claim.isAdmin = true
	}
}

// ReturnToken returns token
func ReturnToken(ctx context.Context, w http.ResponseWriter, r *http.Request, userID, issuer string, options ...TokenClaim) {
	claims := TokenClaims{}
	for _, o := range options {
		o(&claims)
	}
	if claims.isAdmin {
		apicore.ReturnError(ctx, w, r, validation.NewBadRequestError(errors.New("issuing admin token is not implemented yet")))
		return
	}
	token, err := token4auth.IssueAuthToken(ctx, userID, issuer)
	if err != nil {
		err = fmt.Errorf("failed to issue Firebase token: %w", err)
		apicore.ReturnError(ctx, w, r, err)
		return
	}
	header := w.Header()
	//header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Content-Type", "application/json")
	_, _ = w.Write([]byte("{"))
	_, _ = w.Write([]byte(`"token":"`))
	_, _ = w.Write([]byte(token))
	_, _ = w.Write([]byte(`"}`))
}
