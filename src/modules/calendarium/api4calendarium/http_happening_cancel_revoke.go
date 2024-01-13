package api4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/facade4calendarium"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var revokeHappeningCancellation = facade4calendarium.RevokeHappeningCancellation

// httpRevokeHappeningCancellation marks happening as canceled
func httpRevokeHappeningCancellation(w http.ResponseWriter, r *http.Request) {
	var happeningRequest = getHappeningRequestParamsFromURL(r)
	request := dto4calendarium.CancelHappeningRequest{
		HappeningRequest: happeningRequest,
	}
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = revokeHappeningCancellation(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, nil)
}
