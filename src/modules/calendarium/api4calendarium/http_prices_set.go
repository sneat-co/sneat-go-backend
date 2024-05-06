package api4calendarium

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/facade4calendarium"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

func httpSetHappeningPrices(w http.ResponseWriter, r *http.Request) {
	var request dto4calendarium.HappeningPricesRequest
	request.HappeningRequest = getHappeningRequestParamsFromURL(r)
	ctx, userContext, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = facade4calendarium.SetHappeningPrices(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusNoContent, err, nil)
}
