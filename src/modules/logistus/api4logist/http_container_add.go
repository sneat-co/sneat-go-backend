package api4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var addContainers = facade4logist.AddContainers

func httpAddContainers(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.AddContainersRequest
	handler := func(ctx context.Context, userCtx facade.User) (interface{}, error) {
		return nil, addContainers(ctx, userCtx, request)
	}
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, handler, http.StatusNoContent, defaultJsonWithAuthRequired)
}
