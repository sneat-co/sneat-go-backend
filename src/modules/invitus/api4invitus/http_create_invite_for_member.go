package api4invitus

import (
	"context"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"net/http"
)

var createOrReuseInviteForMember = facade4invitus.CreateOrReuseInviteForMember

// httpPostCreateOrReuseInviteForMember supports both POST & GET methods
func httpPostCreateOrReuseInviteForMember(w http.ResponseWriter, r *http.Request) {
	var request facade4invitus.InviteMemberRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			if request.To.Channel == "link" {
				return nil, fmt.Errorf("%w: link invites should be requested via GET", facade.ErrBadRequest)
			}
			request.RemoteClient = apicore.GetRemoteClientInfo(r)
			return createOrReuseInviteForMember(ctx, userCtx, request)
		})
}

// httpGetOrCreateInviteLink gets or creates an invitation link
func httpGetOrCreateInviteLink(w http.ResponseWriter, r *http.Request) {
	var request facade4invitus.InviteMemberRequest
	q := r.URL.Query()
	request.SpaceID = q.Get("space")
	request.To.MemberID = q.Get("member")
	request.To.Channel = "link"
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		httpserver.HandleError(ctx, err, "VerifyRequestAndCreateUserContext", w, r)
		return
	}
	request.RemoteClient = apicore.GetRemoteClientInfo(r)
	response, err := createOrReuseInviteForMember(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
