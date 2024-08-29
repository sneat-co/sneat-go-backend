package api4listus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dto4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/facade4listus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var createList = facade4listus.CreateList

// httpPostCreateList creates a list
func httpPostCreateList(w http.ResponseWriter, r *http.Request) {
	var request dto4listus.CreateListRequest
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := createList(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, &response)
}
