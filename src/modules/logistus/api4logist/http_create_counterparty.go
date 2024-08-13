package api4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var createCounterparty = facade4logist.CreateCounterparty

func httpCreateCounterparty(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.CreateCounterpartyRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx context.Context, userCtx facade.UserContext) (interface{}, error) {
			return createCounterparty(ctx, userCtx, request)
		})
}
