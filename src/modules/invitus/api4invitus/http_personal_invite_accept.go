package api4invitus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/invitus/facade4invitus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

// httpPostAcceptPersonalInvite is an API endpoint that marks a personal invite as accepted
func httpPostAcceptPersonalInvite(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	request := facade4invitus.AcceptPersonalInviteRequest{}
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	request.RemoteClient = apicore.GetRemoteClientInfo(r)
	if err = facade4invitus.AcceptPersonalInvite(ctx, userContext, request); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	apicore.IfNoErrorReturnOK(ctx, w, r, err)
}
