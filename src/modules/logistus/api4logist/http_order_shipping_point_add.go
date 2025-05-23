package api4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var addOrderShippingPoint = facade4logist.AddOrderShippingPoint

func httpAddOrderShippingPoint(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.AddOrderShippingPointRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusOK,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			return addOrderShippingPoint(ctx, request)
		})
}
