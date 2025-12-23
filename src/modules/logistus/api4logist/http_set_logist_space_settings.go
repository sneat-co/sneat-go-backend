package api4logist

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/dto4logist"
	"github.com/sneat-co/sneat-go-backend/src/modules/logistus/facade4logist"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/facade"
)

func httpSetLogistSpaceSettings(w http.ResponseWriter, r *http.Request) {
	var request dto4logist.SetLogistSpaceSettingsRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, defaultJsonWithAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			return nil, facade4logist.SetLogistSpaceSettings(ctx, request)
		})
}
