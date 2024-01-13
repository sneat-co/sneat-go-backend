package api4retrospectus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/meetingus/api4meetingus"
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

// RegisterHttpRoutes registers retrospective routes
func RegisterHttpRoutes(handle modules.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/retrospective/toggle_meeting_timer", api4meetingus.ToggleMeetingTimer(meetingParams))
	handle(http.MethodPost, "/v0/retrospective/toggle_member_timer", api4meetingus.ToggleMemberTimer(meetingParams))

	handle(http.MethodPost, "/v0/retrospective/start_retrospective", httpPostStartRetrospective)
	handle(http.MethodPost, "/v0/retrospective/start_retro_review", httpPostStartRetroReview)

	handle(http.MethodPost, "/v0/retrospective/vote_item", httpPostVoteItem)

	handle(http.MethodPost, "/v0/retrospective/add_retro_item", httpPostAddRetroItem)
	handle(http.MethodPost, "/v0/retrospective/delete_retro_item", httpPostDeleteRetroItem)
	handle(http.MethodPost, "/v0/retrospective/move_retro_item", httpPostMoveRetroItem)
	handle(http.MethodPost, "/v0/retrospective/fix_counts", httpPostFixCounts)
}
