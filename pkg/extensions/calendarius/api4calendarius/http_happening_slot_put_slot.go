package api4calendarius

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/dto4calendarius"
	"github.com/sneat-co/sneat-go-backend/pkg/extensions/calendarius/facade4calendarius"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

func httpAddSlot(w http.ResponseWriter, r *http.Request) {
	httpPutSlot(w, r, facade4calendarius.AddSlot)
}

func httpUpdateSlot(w http.ResponseWriter, r *http.Request) {
	httpPutSlot(w, r, facade4calendarius.UpdateSlot)
}

func httpPutSlot(w http.ResponseWriter, r *http.Request, putMode facade4calendarius.PutMode) {
	var request dto4calendarius.HappeningSlotRequest
	request.HappeningRequest = getHappeningRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = facade4calendarius.UpdateHappeningSlot(ctx, putMode, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusNoContent, err, nil)
}
