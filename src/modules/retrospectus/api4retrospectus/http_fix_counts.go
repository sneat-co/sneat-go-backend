package api4retrospectus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/retrospectus/facade4retrospectus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var fixCounts = facade4retrospectus.FixCounts

// httpPostFixCounts is an API endpoint that triggers fixing of counters in a retrospective
func httpPostFixCounts(w http.ResponseWriter, r *http.Request) {
	ctx, err := verifyRequest(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	request := facade4retrospectus.FixCountsRequest{}
	if err := apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = fixCounts(ctx, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
