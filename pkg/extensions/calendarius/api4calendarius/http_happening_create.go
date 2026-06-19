package api4calendarius

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/facade4calendarius"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

// httpPostCreateHappening creates recurring happening
func httpPostCreateHappening(w http.ResponseWriter, r *http.Request) {
	var request dto4calendarius.CreateHappeningRequest
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	response, err := facade4calendarius.CreateHappening(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, &response)
}
