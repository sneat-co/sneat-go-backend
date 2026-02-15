package api4assetus

import (
	"fmt"
	"net/http"

	"github.com/sneat-co/sneat-core-modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/coretypes"
	"github.com/sneat-co/sneat-go-core/facade"
)

func httpDeleteAsset(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var request dto4spaceus.SpaceItemRequest
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
