package api4contactus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var removeMember = facade4contactus.RemoveMember

// httpPostRemoveMember is an API endpoint that removes a members from a team
func httpPostRemoveMember(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request dto4contactus.ContactRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = removeMember(ctx, userContext, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
