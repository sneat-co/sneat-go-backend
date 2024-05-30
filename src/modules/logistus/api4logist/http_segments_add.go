package api4logist

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpAddSegments(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.AddSegmentsRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			return nil, facade4logist.AddSegments(ctx, userCtx, request)
		})
}
