package api4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var getScrum = facade4scrumus.GetScrum

// httpGetScrum is an API endpoint that returns scrum data
func httpGetScrum(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.Request(verify.AuthenticationRequired(true)))
	if err != nil {
		return
	}
	response, err := getScrum(ctx, facade.IDRequest{ID: r.Header.Get("id")})
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
