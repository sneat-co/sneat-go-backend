package api4calendarium

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/facade4calendarium"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

func httpDeleteHappening(w http.ResponseWriter, r *http.Request) {
	var request = getHappeningRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = facade4calendarium.DeleteHappening(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}
