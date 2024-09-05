package api4auth

import (
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"net/http"
)

func disconnect(w http.ResponseWriter, r *http.Request) {
	if !httpserver.AccessControlAllowOrigin(w, r) {
		return
	}

	ctx, userCtx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.Request(verify.AuthenticationRequired(true)))
	if err != nil {
		return
	}
	err = facade4auth.Disconnect(ctx, userCtx, r.URL.Query().Get("provider"))
	apicore.ReturnStatus(ctx, w, r, http.StatusNoContent, err)
}
