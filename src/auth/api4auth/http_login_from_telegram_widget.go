package api4auth

import (
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/sneat-co/sneat-go-backend/src/auth/dto4auth"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"net/http"
)

func httpLoginFromTelegramWidget(w http.ResponseWriter, r *http.Request) {

	ctx, err := apicore.VerifyRequest(w, r, verify.DefaultJsonWithNoAuthRequired)
	if err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	var tgAuthData dto4auth.TelegramAuthData
	if err = apicore.DecodeRequestBody(w, r, &tgAuthData); err != nil {
		// apicore.ReturnError(ctx, w, r, err)
		return
	}

	var initData twainitdata.InitData
	initData.Hash = tgAuthData.Hash
	initData.AuthDateRaw = tgAuthData.AuthDate
	initData.User.ID = tgAuthData.ID
	initData.User.Username = tgAuthData.Username
	initData.User.FirstName = tgAuthData.FirstName
	initData.User.LastName = tgAuthData.LastName
	initData.User.PhotoURL = tgAuthData.PhotoURL

	signInWithTelegram(ctx, w, r, initData)
}
