package api4teamus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var joinTeam = facade4invitus.JoinTeam

// httpPostJoinTeam joins a members to a team
func httpPostJoinTeam(w http.ResponseWriter, r *http.Request) {
	var request facade4invitus.JoinTeamRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, func(ctx context.Context, userCtx facade.User) (response interface{}, err error) {
		return joinTeam(ctx, userCtx, request)
	}, http.StatusOK, verify.DefaultJsonWithAuthRequired)
}
