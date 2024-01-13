package api4retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/facade4retrospectus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var deleteRetroItem = facade4retrospectus.DeleteRetroItem

// httpPostDeleteRetroItem is an API endpoint that removes an items from a retrospective
func httpPostDeleteRetroItem(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := verifyAuthorizedJSONRequest(w, r, verify.MinJSONRequestSize, 10*verify.KB)
	if err != nil {
		return
	}
	request := facade4retrospectus.RetroItemRequest{}
	if err := apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = deleteRetroItem(ctx, userContext, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
