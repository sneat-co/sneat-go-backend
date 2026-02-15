package api4logist

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
)

var createOrder = facade4logist.CreateOrder

func httpCreateOrder(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.CreateOrderRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			order, err := createOrder(ctx, request)
			response := dto4logist.CreateOrderResponse{
				Order: order,
			}
			return response, err
		})
}
