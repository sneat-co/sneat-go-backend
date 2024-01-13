package api4userus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dto4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/models4userus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

var initUserRecord = facade4userus.InitUserRecord

// httpInitUserRecord sets user title
func httpInitUserRecord(w http.ResponseWriter, r *http.Request) {
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request dto4userus.InitUserRecordRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	request.RemoteClient = apicore.GetRemoteClientInfo(r)
	var user models4userus.UserContext
	user, err = initUserRecord(ctx, userContext, request)
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, user.Dto)
}
