package api4invitus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-backend/src/coremodules/spaceus/dto4spaceus"
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
		SpaceRequest: dto4spaceus.SpaceRequest{
			SpaceID: strings.TrimSpace(q.Get("spaceID")),
		},
		InviteID: strings.TrimSpace(q.Get("inviteID")),
	}
	response, err := facade4invitus.GetPersonal(ctx, user, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
