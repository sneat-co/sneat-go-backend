package api4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpSetContainerPointFreightFields(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.SetContainerPointFreightFieldsRequest
	handler := func(ctx context.Context, userCtx facade.User) (interface{}, error) {
		return nil, facade4logist.SetContainerPointFreightFields(ctx, userCtx, request)
	}
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, handler, http.StatusNoContent, defaultJsonWithAuthRequired)
}