package api4teamus

import (
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

// RegisterHttpRoutes registers team routes
func RegisterHttpRoutes(handle modules.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/teams/create_team", httpPostCreateTeam)
	//
	handle(http.MethodPost, "/v0/team/join_info", httpPostGetTeamJoinInfo)
	handle(http.MethodPost, "/v0/team/join_team", httpPostJoinTeam)
	handle(http.MethodPost, "/v0/team/refuse_to_join_team", httpPostRefuseToJoinTeam)
	handle(http.MethodPost, "/v0/team/add_metric", httpPostAddMetric)
	handle(http.MethodPost, "/v0/team/remove_metrics", httpPostRemoveMetrics)
	handle(http.MethodGet, "/v0/team", httpGetTeam)
}
