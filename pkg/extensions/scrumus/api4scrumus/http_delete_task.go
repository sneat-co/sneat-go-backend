package api4scrumus

import (
	"net/http"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/meetingus/facade4meetingus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/scrumus/facade4scrumus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/httpserver"
)

var deleteTask = facade4scrumus.DeleteTask

// httpDeleteTask is an API endpoint that delete a task
func httpDeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		httpserver.HandleError(ctx, err, "httpDeleteTask", w, r)
		return
	}
	query := r.URL.Query()
	request := facade4scrumus.DeleteTaskRequest{
		Request: facade4meetingus.Request{
			SpaceRequest: dto4spaceus.SpaceRequest{
				SpaceID: coretypes.SpaceID(query.Get("space")),
			},
			MeetingID: query.Get("date"),
		},
		Task:      query.Get("id"),
		Type:      query.Get("type"),
		ContactID: query.Get("members"),
	}

	if err := request.Validate(); err != nil {
		httpserver.HandleError(ctx, err, "httpDeleteTask", w, r)
		return
	}
	err = deleteTask(ctx, request)
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
