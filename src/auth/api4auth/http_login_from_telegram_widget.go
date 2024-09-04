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
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	initData := twainitdata.InitData{
		Hash:        tgAuthData.Hash,
		AuthDateRaw: tgAuthData.AuthDate,
		User: twainitdata.User{
			ID:        tgAuthData.ID,
			Username:  tgAuthData.Username,
			FirstName: tgAuthData.FirstName,
			LastName:  tgAuthData.LastName,
			PhotoURL:  tgAuthData.PhotoURL,
		},
	}

	signInWithTelegram(ctx, w, r, "", initData)
}
