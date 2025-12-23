package api4listus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

// httpDeleteListItems deletes list items
func httpDeleteListItems(w http.ResponseWriter, r *http.Request) {
	var request dto4listus.ListItemIDsRequest
	request.ListRequest = getListRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	if err = request.Validate(); err != nil {
		apicore.ReturnError(r.Context(), w, r, err)
		return
	}
	_, _, err = facade4listus.DeleteListItems(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}
