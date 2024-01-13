package api4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var createListItems = facade4listus.CreateListItems

// httpPostCreateListItems creates list items
func httpPostCreateListItems(w http.ResponseWriter, r *http.Request) {
	var request facade4listus.CreateListItemsRequest
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := createListItems(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, &response)
}
