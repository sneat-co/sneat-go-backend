package api4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var setListItemsIsDone = facade4listus.SetListItemsIsDone

// httpPostSetListItemsIsDone marks list items as completed
func httpPostSetListItemsIsDone(w http.ResponseWriter, r *http.Request) {
	var request dto4listus.ListItemsSetIsDoneRequest
	request.ListRequest = getListRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	if err = request.Validate(); err != nil {
		apicore.ReturnError(r.Context(), w, r, err)
		return
	}
	_, _, err = setListItemsIsDone(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusNoContent, err, nil)
}
