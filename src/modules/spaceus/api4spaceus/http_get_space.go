package api4spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4spaceus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

//var getSpaceByID = facade4spaceus.GetSpaceByID

// httpGetSpace is an API endpoint that return team data
func httpGetSpace(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	verifyOptions := verify.Request(verify.AuthenticationRequired(true))
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verifyOptions)
	if err != nil {
		return
	}
	var space dbo4spaceus.SpaceEntry
	var response any
	if space, err = facade4spaceus.GetSpace(ctx, userContext, id); err == nil {
		response = space.Data
	}
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
