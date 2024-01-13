package api4scrumus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/api4meetingus"
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

// RegisterHttpRoutes registers scrum routes
func RegisterHttpRoutes(handle modules.HTTPHandleFunc) {
	handle(http.MethodGet, "/v0/scrum", httpGetScrum)
	handle(http.MethodPost, "/v0/scrum/add_task", httpPostAddTask)
	handle(http.MethodPost, "/v0/scrum/set_metric", httpPostSetMetric)
	handle(http.MethodPost, "/v0/scrum/reorder_task", httpPostReorderTask)
	handle(http.MethodPost, "/v0/scrum/add_comment", httpPostAddComment)
	handle(http.MethodDelete, "/v0/scrum/delete_task", httpDeleteTask)
	handle(http.MethodPost, "/v0/scrum/thumb_up_task", httpPostThumbUp)
	handle(http.MethodPost, "/v0/scrum/toggle_meeting_timer", api4meetingus.ToggleMeetingTimer(meetingParams))
	handle(http.MethodPost, "/v0/scrum/toggle_member_timer", api4meetingus.ToggleMemberTimer(meetingParams))
}
