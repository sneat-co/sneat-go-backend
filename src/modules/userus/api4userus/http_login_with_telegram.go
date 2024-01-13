package api4userus

import (
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dto4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

func httpLoginWithTelegram(w http.ResponseWriter, r *http.Request) {

	ctx, err := apicore.VerifyRequest(w, r, verify.DefaultJsonWithNoAuthRequired)
	if err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	var tgAuthData dto4userus.TelegramAuthData
	if err = apicore.DecodeRequestBody(w, r, &tgAuthData); err != nil {
		return
	}

	var response facade4userus.FirebaseCustomAuthResponse
	response, err = facade4userus.LoginWithTelegram(ctx, tgAuthData)

	apicore.ReturnJSON(ctx, w, r, http.StatusOK, err, response)
}
