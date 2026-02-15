package api4logist

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
)

var setOrderCounterparties = facade4logist.SetOrderCounterparties

func httpSetOrderCounterparties(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.SetOrderCounterpartiesRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusOK,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			return setOrderCounterparties(ctx, request)
		})
}
