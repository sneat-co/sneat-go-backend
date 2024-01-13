package api4invitus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
	"strings"
)

// httpGetPersonal is an API endpoint that returns personal invite data
func httpGetPersonal(w http.ResponseWriter, r *http.Request) {
	ctx, user, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	q := r.URL.Query()
	request := facade4invitus.GetPersonalInviteRequest{
		TeamRequest: dto4teamus.TeamRequest{
			TeamID: strings.TrimSpace(q.Get("teamID")),
		},
		InviteID: strings.TrimSpace(q.Get("inviteID")),
	}
	response, err := facade4invitus.GetPersonal(ctx, user, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
