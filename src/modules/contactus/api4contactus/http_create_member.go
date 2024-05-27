package api4contactus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dal4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var createMember = facade4contactus.CreateMember

// httpPostCreateMember is an API endpoint that adds a members to a team.
// While is very similar to contactus/api4contactus/http_create_contact.go, it's not the same.
func httpPostCreateMember(w http.ResponseWriter, r *http.Request) {
	var request dal4contactus.CreateMemberRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			return createMember(ctx, userCtx, request)
		})
}
