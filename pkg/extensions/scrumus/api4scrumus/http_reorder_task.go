package api4scrumus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var reorderTask = facade4scrumus.ReorderTask

// httpPostReorderTask is an API endpoints that reorders tasks
func httpPostReorderTask(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request facade4scrumus.ReorderTaskRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	err = reorderTask(ctx, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
