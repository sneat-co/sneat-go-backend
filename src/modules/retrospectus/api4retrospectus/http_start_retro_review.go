package api4retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/facade4retrospectus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var startRetroReview = facade4retrospectus.StartRetroReview

// httpPostStartRetroReview is an API endpoint that starts retrospective
func httpPostStartRetroReview(w http.ResponseWriter, r *http.Request) {
	ctx, err := verifyRequest(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	request := facade4retrospectus.RetroRequest{}
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	response, err := startRetroReview(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
