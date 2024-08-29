package api4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var deleteListItems = facade4listus.DeleteListItems

// httpDeleteListItems deletes list items
func httpDeleteListItems(w http.ResponseWriter, r *http.Request) {
	var request dto4listus.ListItemIDsRequest
	request.ListRequest = getListRequestParamsFromURL(r)
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	if err = request.Validate(); err != nil {
		apicore.ReturnError(r.Context(), w, r, err)
		return
	}
	err = deleteListItems(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}
