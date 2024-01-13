package api4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/facade4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var getTeam = facade4teamus.GetTeam

//var getTeamByID = facade4teamus.GetTeamByID

// httpGetTeam is an API endpoint that return team data
func httpGetTeam(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	verifyOptions := verify.Request(verify.AuthenticationRequired(true))
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verifyOptions)
	if err != nil {
		return
	}
	var team dal4teamus.TeamContext
	team, err = getTeam(ctx, userContext, id)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, team.Data)
}
