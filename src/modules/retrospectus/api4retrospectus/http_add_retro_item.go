package api4retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/facade4retrospectus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var addRetroItem = facade4retrospectus.AddRetroItem

// httpPostAddRetroItem adds an item to a retrospective
func httpPostAddRetroItem(w http.ResponseWriter, r *http.Request) {
	ctx, err := verifyAuthorizedJSONRequest(w, r, verify.MinJSONRequestSize, 10*verify.KB)
	if err != nil {
		return
	}
	request := facade4retrospectus.AddRetroItemRequest{}
	if err := apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	response, err := addRetroItem(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, response)
}
