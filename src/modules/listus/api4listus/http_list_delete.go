package api4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var deleteList = facade4listus.DeleteList

// httpDeleteList deletes a list
func httpDeleteList(w http.ResponseWriter, r *http.Request) {
	var request = getListRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = deleteList(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}
