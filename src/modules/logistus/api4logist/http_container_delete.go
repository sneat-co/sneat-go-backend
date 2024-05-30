package api4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

var deleteContainer = facade4logist.DeleteContainer

func httpDeleteContainer(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.ContainerRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			return nil, deleteContainer(ctx, userCtx, request)
		})
}
