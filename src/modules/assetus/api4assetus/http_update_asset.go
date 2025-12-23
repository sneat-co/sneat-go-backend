package api4assetus

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

func httpPostUpdateAsset(w http.ResponseWriter, r *http.Request) {
	var request dto4assetus.UpdateAssetRequest
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			return nil, facade4assetus.UpdateAsset(ctx, request)
		})
}
