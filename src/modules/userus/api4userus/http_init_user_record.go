package api4userus

import (
	"errors"
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dto4userus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

// httpInitUserRecord sets user title
func httpInitUserRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != "" {
		apicore.ReturnError(r.Context(), w, r, errors.New("temporary disabled"))
		return
	}
	ctx, userContext, err := apicore.VerifyRequestAndCreateUserContext(w, r, verify.DefaultJsonWithAuthRequired)
	if err != nil {
		return
	}
	var request dto4userus.InitUserRecordRequest
	if err = apicore.DecodeRequestBody(w, r, &request); err != nil {
		return
	}
	request.RemoteClient = apicore.GetRemoteClientInfo(r)
	var user dbo4userus.UserEntry
	userToCreate := facade4auth.DataToCreateUser{
		AuthProvider:    request.AuthProvider,
		Email:           request.Email,
		EmailIsVerified: request.EmailIsVerified,
		IanaTimezone:    request.IanaTimezone,
		RemoteClient:    request.RemoteClient,
	}
	if request.Names != nil {
		userToCreate.Names = *request.Names
	}
	if user, err = facade4auth.CreateUserRecords(ctx, userContext, userToCreate); err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}
	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, user.Data)
}
