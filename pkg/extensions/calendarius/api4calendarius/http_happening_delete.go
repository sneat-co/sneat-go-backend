package api4calendarius

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/facade4calendarius"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

func httpDeleteHappening(w http.ResponseWriter, r *http.Request) {
	var request = getHappeningRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = facade4calendarius.DeleteHappening(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}
