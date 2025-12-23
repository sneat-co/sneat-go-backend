package api4calendarium

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/facade4calendarium"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

func httpSetHappeningPrices(w http.ResponseWriter, r *http.Request) {
	var request dto4calendarium.HappeningPricesRequest
	request.HappeningRequest = getHappeningRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = facade4calendarium.SetHappeningPrices(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusNoContent, err, nil)
}
