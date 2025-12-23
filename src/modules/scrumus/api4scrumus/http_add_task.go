package api4scrumus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

var addTask = facade4scrumus.AddTask

// httpPostAddTask adds a task
func httpPostAddTask(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request facade4scrumus.AddTaskRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	response, err := addTask(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, response)
}
