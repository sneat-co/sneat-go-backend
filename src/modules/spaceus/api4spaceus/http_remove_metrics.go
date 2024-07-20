package api4spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var removeMetrics = facade4teamus.RemoveMetrics

// httpPostRemoveMetrics is an API endpoint that removes a team metric
func httpPostRemoveMetrics(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request dto4spaceus.SpaceMetricsRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = removeMetrics(ctx, userContext, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
