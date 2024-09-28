package api4spaceus

import (
	"github.com/sneat-co/sneat-go-core/module"
	"net/http"
)

// RegisterHttpRoutes registers team routes
func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/spaces/create_space", httpPostCreateSpace)
	//
	handle(http.MethodPost, "/v0/space/join_info", httpPostGetSpaceJoinInfo)
	handle(http.MethodPost, "/v0/space/join_team", httpPostJoinSpace)
	handle(http.MethodPost, "/v0/space/refuse_to_join_team", httpPostRefuseToJoinSpace)
	handle(http.MethodPost, "/v0/space/add_metric", httpPostAddMetric)
	handle(http.MethodPost, "/v0/space/remove_metrics", httpPostRemoveMetrics)
	handle(http.MethodGet, "/v0/space", httpGetSpace)
}
