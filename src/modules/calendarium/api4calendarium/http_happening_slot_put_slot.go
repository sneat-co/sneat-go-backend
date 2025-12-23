package api4calendarium

import (
	"net/http"

	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/dto4calendarium"
	"github.com/sneat-co/sneat-go-backend/src/modules/calendarium/facade4calendarium"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
)

func httpAddSlot(w http.ResponseWriter, r *http.Request) {
	httpPutSlot(w, r, facade4calendarium.AddSlot)
}

func httpUpdateSlot(w http.ResponseWriter, r *http.Request) {
	httpPutSlot(w, r, facade4calendarium.UpdateSlot)
}

func httpPutSlot(w http.ResponseWriter, r *http.Request, putMode facade4calendarium.PutMode) {
	var request dto4calendarium.HappeningSlotRequest
	request.HappeningRequest = getHappeningRequestParamsFromURL(r)
	ctx, err := apicore.VerifyAuthenticatedRequestAndDecodeBody(w, r, verify.DefaultJsonWithAuthRequired, &request)
	if err != nil {
		return
	}
	err = facade4calendarium.UpdateHappeningSlot(ctx, putMode, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusNoContent, err, nil)
}
