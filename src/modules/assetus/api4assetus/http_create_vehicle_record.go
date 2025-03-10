package api4assetus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/dto4assetus"
	"github.com/sneat-co/sneat-go-backend/src/modules/assetus/facade4assetus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/sneat-co/sneat-go-core/facade"
	"net/http"
)

func httpPostCreateVehicleRecord(w http.ResponseWriter, r *http.Request) {
	var (
		request dto4assetus.AddVehicleRecordRequest
	)
	apicore.HandleAuthenticatedRequestWithBody(w, r, &request, verify.DefaultJsonWithAuthRequired, http.StatusOK,
		func(ctx facade.ContextWithUser) (interface{}, error) {
			return facade4assetus.AddVehicleRecord(ctx, request)
		},
	)
}
