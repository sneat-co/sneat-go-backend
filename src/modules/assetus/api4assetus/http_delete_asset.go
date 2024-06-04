package api4assetus

import (
	"context"
	"fmt"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dal4teamus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

// httpDeleteAsset deletes assets
func httpDeleteAsset(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var request dal4teamus.TeamItemRequest
	request.TeamID = q.Get("team")
	request.ID = q.Get("id")
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.NoContentAuthRequired, http.StatusNoContent,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			if err := facade4assetus.DeleteAsset(ctx, userCtx, request); err != nil {
				return nil, fmt.Errorf("failed to delete asset: %w", err)
			}
			return nil, nil
		})
}
