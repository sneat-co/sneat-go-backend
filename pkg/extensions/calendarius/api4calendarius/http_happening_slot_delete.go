package api4calendarius

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/facade4calendarius"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

func httpDeleteSlot(w http.ResponseWriter, r *http.Request) {
	var request = dto4calendarius.DeleteHappeningSlotRequest{
		HappeningSlotRefRequest: dto4calendarius.HappeningSlotRefRequest{
			HappeningRequest: getHappeningRequestParamsFromURL(r),
		},
	}
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = facade4calendarius.DeleteSlot(ctx, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusCreated, err, nil)
}
