package api4logist

import (
	"github.com/sneat-co/sneat-go-core/module"
	"net/http"
)

const (
	RoutePathDeleteOrderCounterparty = "/v0/logistus/order/delete_order_counterparty"
	RoutePathSetOrderCounterparties  = "/v0/logistus/order/set_order_counterparties"
	RoutePathSetOrderStatus          = "/v0/logistus/order/set_order_status"
	RoutePathCreateOrder             = "/v0/logistus/create_order"
	RoutePathSetLogistSpaceSettings  = "/v0/logistus/set_logist_team_settings"
	RoutePathCreateCounterparty      = "/v0/logistus/create_counterparty"
	RoutePathOrderDeleteContainer    = "/v0/logistus/order/delete_container"

	RoutePathAddOrderShippingPoint = "/v0/logistus/order/add_shipping_point"
	RoutePathDeleteShippingPoint   = "/v0/logistus/order/delete_shipping_point"
	RoutePathUpdateShippingPoint   = "/v0/logistus/order/update_shipping_point"
	RoutePathAddContainers         = "/v0/logistus/order/add_containers"
	RoutePathSetContainerFields    = "/v0/logistus/order/set_container_fields"

	RoutePathAddSegments                    = "/v0/logistus/order/add_segments"
	RoutePathDeleteSegments                 = "/v0/logistus/order/delete_segments"
	RoutePathAddContainerPoints             = "/v0/logistus/order/add_container_points"
	RoutePathUpdateContainerPoint           = "/v0/logistus/order/update_container_point"
	RoutePathSetContainerPointTask          = "/v0/logistus/order/set_container_point_task"
	RoutePathSetContainerPointFields        = "/v0/logistus/order/set_container_point_fields"
	RoutePathSetContainerPointFreightFields = "/v0/logistus/order/set_container_point_freight_fields"
	RoutePathSetContainerEndpointFields     = "/v0/logistus/order/set_container_endpoint_fields"
	RoutePathDeleteContainerPoints          = "/v0/logistus/order/delete_container_points"
)

// RegisterHttpRoutes registers logistus routes
func RegisterHttpRoutes(handle module.HTTPHandleFunc) {
	handle(http.MethodPost, RoutePathSetLogistSpaceSettings, httpSetLogistSpaceSettings)
	handle(http.MethodPost, RoutePathCreateCounterparty, httpCreateCounterparty)
	handle(http.MethodPost, RoutePathCreateOrder, httpCreateOrder)
	handle(http.MethodPost, RoutePathSetOrderStatus, httpSetOrderStatus)
	handle(http.MethodPost, RoutePathSetOrderCounterparties, httpSetOrderCounterparties)
	handle(http.MethodDelete, RoutePathDeleteOrderCounterparty, httpDeleteOrderCounterparty)
	//
	handle(http.MethodPost, RoutePathAddOrderShippingPoint, httpAddOrderShippingPoint)
	handle(http.MethodPost, RoutePathUpdateShippingPoint, httpUpdateShippingPoint)
	handle(http.MethodDelete, RoutePathDeleteShippingPoint, httpDeleteShippingPoint)
	//
	handle(http.MethodPost, RoutePathAddContainers, httpAddContainers)
	handle(http.MethodDelete, RoutePathOrderDeleteContainer, httpDeleteContainer)
	handle(http.MethodPost, RoutePathSetContainerFields, httpSetContainerFields)
	//
	handle(http.MethodPost, RoutePathAddSegments, httpAddSegments)
	handle(http.MethodDelete, RoutePathDeleteSegments, httpDeleteSegments)
	//
	handle(http.MethodPost, RoutePathAddContainerPoints, httpAddContainerPoints)
	handle(http.MethodPost, RoutePathSetContainerPointTask, httpSetContainerPointTask)
	handle(http.MethodPost, RoutePathSetContainerEndpointFields, httpSetContainerEndpointFields)
	handle(http.MethodPost, RoutePathSetContainerPointFields, httpSetContainerPointFields)
	handle(http.MethodPost, RoutePathSetContainerPointFreightFields, httpSetContainerPointFreightFields)

	handle(http.MethodPost, RoutePathUpdateContainerPoint, httpUpdateContainerPoint)
	handle(http.MethodDelete, RoutePathDeleteContainerPoints, httpDeleteContainerPoints)
}
