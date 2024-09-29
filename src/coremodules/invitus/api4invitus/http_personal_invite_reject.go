package api4invitus

import (
	"github.com/sneat-co/sneat-go-backend/src/coremodules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

// httpPostRejectPersonalInvite rejects personal invite
func httpPostRejectPersonalInvite(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithNoAuthRequired)
	if err != nil {
		return
	}
	request := facade4invitus.RejectPersonalInviteRequest{}
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	if err = facade4invitus.RejectPersonalInvite(ctx, userContext, request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
