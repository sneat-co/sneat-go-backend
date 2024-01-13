package api4contactus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/dto4contactus"
	"github.com/sneat-co/sneat-go-backend/src/modules/contactus/facade4contactus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var setContactsStatus = facade4contactus.SetContactsStatus

func httpSetContactStatus(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.SetContactsStatusRequest
	handler := func(ctx context.Context, userCtx facade.User) (interface{}, error) {
		return nil, setContactsStatus(ctx, userCtx, request)
	}
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, handler, http.StatusCreated, verify.DefaultJsonWithAuthRequired)
}
