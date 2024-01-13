package api4userus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpSetUserCountry(w http.ResponseWriter, r *http.Request) {
	var request facade4userus.SetUserCountryRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, func(ctx context.Context, userCtx facade.User) (response interface{}, err error) {
		return nil, facade4userus.SetUserCountry(ctx, userCtx, request)
	}, http.StatusNoContent, verify.DefaultJsonWithAuthRequired)
}
