package api4teamus

import (
	"github.com/sneat-co/sneat-go-core/modules"
	"net/http"
)

// RegisterHttpRoutes registers team routes
func RegisterHttpRoutes(handle modules.HTTPHandleFunc) {
	handle(http.MethodPost, "/v0/space/create_space", httpPostCreateSpace)
	//
	handle(http.MethodPost, "/v0/space/join_info", httpPostGetSpaceJoinInfo)
	handle(http.MethodPost, "/v0/space/join_team", httpPostJoinSpace)
	handle(http.MethodPost, "/v0/space/refuse_to_join_team", httpPostRefuseToJoinSpace)
	handle(http.MethodPost, "/v0/space/add_metric", httpPostAddMetric)
	handle(http.MethodPost, "/v0/space/remove_metrics", httpPostRemoveMetrics)
	handle(http.MethodGet, "/v0/space", httpGetSpace)
}
