package api4auth

import (
	"context"
	"fmt"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"net/http"
)

func signInWithTelegram(ctx context.Context, w http.ResponseWriter, r *http.Request, botID string, initData twainitdata.InitData) {
	var (
		err       error
		tgBotUser facade4auth.BotUserEntry
	)
	remoteClientInfo := apicore.GetRemoteClientInfo(r)
	if tgBotUser, _, err = facade4auth.SignInWithTelegram(ctx, botID, initData, remoteClientInfo); err != nil {
		apicore.ReturnError(ctx, w, r, fmt.Errorf("failed to sign in with Telegram: %w", err))
		return
	}

	appUserID := tgBotUser.Data.GetAppUserID()
	api4debtus.ReturnToken(ctx, w, appUserID, telegram.PlatformID)
}
