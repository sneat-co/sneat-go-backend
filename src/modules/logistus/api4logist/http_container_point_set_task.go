package api4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpSetContainerPointTask(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.SetContainerPointTaskRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx context.Context, userCtx facade.UserContext) (interface{}, error) {
			return nil, facade4logist.SetContainerPointTask(ctx, userCtx, request)
		})
}
