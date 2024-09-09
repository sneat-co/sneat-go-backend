package api4invitus

import (
	"context"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"github.com/strongo/validation"
	"net/http"
)

var createOrReuseInviteForMember = facade4invitus.CreateOrReuseInviteForMember

// httpPostCreateOrReuseInviteForMember supports both POST & GET methods
func httpPostCreateOrReuseInviteForMember(w http.ResponseWriter, r *http.Request) {
	var request facade4invitus.InviteMemberRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx context.Context, userCtx facade.UserContext) (interface{}, error) {
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

	if request.SpaceID = q.Get("space"); request.SpaceID == "" {
		apicore.ReturnError(r.Context(), w, r, validation.NewErrRequestIsMissingRequiredField("space"))
		// TODO(deprecate): httpserver.HandleError(nil, validation.NewErrRequestIsMissingRequiredField("space"), "httpGetOrCreateInviteLink", w, r)
		return
	}
	if request.To.MemberID = q.Get("member"); request.To.MemberID == "" {
		apicore.ReturnError(r.Context(), w, r, validation.NewErrRequestIsMissingRequiredField("member"))
		return
	}

	request.To.Channel = "link"
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.Request(
		verify.AuthenticationRequired(true),
		verify.MaximumContentLength(0),
	))
	if err != nil {
		httpserver.HandleError(ctx, err, "VerifyRequestAndCreateUserContext", w, r)
		return
	}
	request.RemoteClient = apicore.GetRemoteClientInfo(r)
	response, err := createOrReuseInviteForMember(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
