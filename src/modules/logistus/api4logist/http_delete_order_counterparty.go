package api4logist

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var deleteOrderCounterparty = facade4logist.DeleteOrderCounterparty

func httpDeleteOrderCounterparty(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.DeleteOrderCounterpartyRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			return nil, deleteOrderCounterparty(ctx, request)
		})
}
