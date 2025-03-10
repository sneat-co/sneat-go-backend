package api4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var setMetric = facade4scrumus.SetMetric

// httpPostSetMetric is an API endpoint that sets metric value
func httpPostSetMetric(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request facade4scrumus.SetMetricRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	response, err := setMetric(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, response)
}
