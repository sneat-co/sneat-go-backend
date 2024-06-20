package api4assetus

import (
	"context"
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

func httpPostCreateVehicleRecord(w http.ResponseWriter, r *http.Request) {
	var (
		request dto4assetus.CreateAssetRequest
	)
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusCreated,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			// asset, err := facade4assetus.CreateAsset(ctx, userCtx, request)
			// if err != nil {
			// 	return nil, fmt.Errorf("failed to create asset: %w", err)
			// }
			// if asset.ID == "" {
			// 	return nil, errors.New("asset created by facade does not have an ContactID")
			// }
			// if err = asset.Data.Validate(); err != nil {
			// 	err = fmt.Errorf("asset created by facade is not valid: %w", err)
			// 	return asset, err
			// }
			// return asset, nil
		},
	)
}
