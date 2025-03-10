package api4assetus

import (
	"fmt"
	"github.com/sneat-co/sneat-core-modules/spaceus/dal4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpDeleteAsset(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var request dal4spaceus.SpaceItemRequest
	request.SpaceID = coretypes.SpaceID(q.Get("space"))
	request.ID = q.Get("id")
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.NoContentAuthRequired, http.StatusNoContent,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			if err := facade4assetus.DeleteAsset(ctx, request); err != nil {
				return nil, fmt.Errorf("failed to delete asset: %w", err)
			}
			return nil, nil
		})
}
