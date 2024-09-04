package api4auth

import (
	"github.com/bots-go-framework/bots-fw-telegram-webapp/tgloginwidget"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/apicore/verify"
	"github.com/strongo/validation"
	"net/http"
	"os"
	"strings"
)

func httpLoginFromTelegramWidget(w http.ResponseWriter, r *http.Request) {

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

	botID := os.Getenv("TELEGRAM_BOT_ID")
	if botID == "" {
		botID = "SneatBot"
	}
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN_" + strings.ToUpper(botID))

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

	signInWithTelegram(ctx, w, r, "", initData)
}
