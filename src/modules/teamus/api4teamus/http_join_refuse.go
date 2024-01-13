package api4teamus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
	"strconv"
)

var refuseToJoinTeam = facade4contactus.RefuseToJoinTeam

// httpPostRefuseToJoinTeam an API endpoint that records user refusal to join a team
func httpPostRefuseToJoinTeam(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithNoAuthRequired)
	if err != nil {
		return
	}
	q := r.URL.Query()
	var pin int
	if pin, err = strconv.Atoi(q.Get("pin")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("pin is expected to be an integer"))
		return
	}
	request := facade4contactus.RefuseToJoinTeamRequest{
		TeamID: q.Get("id"),
		Pin:    int32(pin),
	}
	err = refuseToJoinTeam(ctx, userContext, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
