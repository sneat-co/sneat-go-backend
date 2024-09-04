package api4auth

import (
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/tgloginwidget"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/sneat-co/sneat-go-backend/src/botscore"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/strongo/validation"
	"net/http"
)

func httpLoginFromTelegramWidget(w http.ResponseWriter, r *http.Request) {

	botID := r.URL.Query().Get("botID")
	if botID == "" {
		apicore.ReturnError(r.Context(), w, r, validation.NewErrRecordIsMissingRequiredField("botID"))
		return
	}

	ctx, err := apicore.VerifyRequest(w, r, verify.DefaultJsonWithNoAuthRequired)
	if err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	var tgAuthData tgloginwidget.AuthData
	if err = apicore.DecodeRequestBody(w, r, &tgAuthData); err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	var botToken string
	if botToken, err = botscore.GetBotToken(telegram.PlatformID, botID); err != nil {
		return
	}

	if err = tgAuthData.Check(botToken); err != nil {
		apicore.ReturnError(ctx, w, r, validation.NewBadRequestError(err))
		return
	}

	initData := twainitdata.InitData{
		Hash:        tgAuthData.Hash,
		AuthDateRaw: int(tgAuthData.AuthDate),
		User: twainitdata.User{
			ID:        tgAuthData.ID,
			Username:  tgAuthData.Username,
			FirstName: tgAuthData.FirstName,
			LastName:  tgAuthData.LastName,
			PhotoURL:  tgAuthData.PhotoURL,
		},
	}

	signInWithTelegram(ctx, w, r, botID, initData)
}
