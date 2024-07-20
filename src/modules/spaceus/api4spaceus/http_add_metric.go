package api4spaceus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/facade4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
	"strings"
)

// httpPostAddMetric is an API endpoint that adds a metric
func httpPostAddMetric(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request facade4teamus.AddSpaceMetricRequest
	if request.SpaceID = r.URL.Query().Get("id"); strings.TrimSpace(request.SpaceID) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("space 'id' should be passed as query parameter"))
		return
	}
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = addMetric(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}

var addMetric = facade4teamus.AddMetric
