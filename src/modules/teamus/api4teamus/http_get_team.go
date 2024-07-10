package api4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/facade4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

//var getSpaceByID = facade4teamus.GetSpaceByID

// httpGetSpace is an API endpoint that return team data
func httpGetSpace(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	verifyOptions := verify.Request(verify.AuthenticationRequired(true))
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verifyOptions)
	if err != nil {
		return
	}
	var space dal4teamus.SpaceEntry
	var response any
	if space, err = facade4teamus.GetSpace(ctx, userContext, id); err == nil {
		response = space.Data
	}
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
