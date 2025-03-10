package api4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var reorderListItem = facade4listus.ReorderListItem

// httpPostReorderListItem reorders list items
func httpPostReorderListItem(w http.ResponseWriter, r *http.Request) {
	var request dto4listus.ReorderListItemsRequest
	request.ListRequest = getListRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	if err = request.Validate(); err != nil {
		apicore.ReturnError(r.Context(), w, r, err)
		return
	}
	err = reorderListItem(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, nil)
}
