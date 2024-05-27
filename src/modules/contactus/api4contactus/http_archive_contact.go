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

// httpPostArchiveContact is an API endpoint that archive contact - e.g., hides it from the list of contacts
func httpPostArchiveContact(w http.ResponseWriter, r *http.Request) {
	var request dto4contactus.ContactRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusOK,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			return nil, facade4contactus.ArchiveContact(ctx, userCtx, request)
		})
}
