package api4assetus

import (
	"context"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/facade4assetus"
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
)

func httpPostCreateVehicleRecord(w http.ResponseWriter, r *http.Request) {
	var (
		request dto4assetus.AddVehicleRecordRequest
	)
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusOK,
		func(ctx context.Context, userCtx facade.User) (interface{}, error) {
			return facade4assetus.AddVehicleRecord(ctx, userCtx, request)
		},
	)
}
