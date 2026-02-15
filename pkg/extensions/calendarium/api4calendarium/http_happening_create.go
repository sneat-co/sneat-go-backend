package api4calendarium

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarium/facade4calendarium"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

// httpPostCreateHappening creates recurring happening
func httpPostCreateHappening(w http.ResponseWriter, r *http.Request) {
	var request dto4calendarium.CreateHappeningRequest
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := facade4calendarium.CreateHappening(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, &response)
}
