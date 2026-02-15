package api4logist

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
)

var addContainerPoints = facade4logist.AddContainerPoints

func httpAddContainerPoints(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.AddContainerPointsRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			return nil, addContainerPoints(ctx, request)
		})
}
