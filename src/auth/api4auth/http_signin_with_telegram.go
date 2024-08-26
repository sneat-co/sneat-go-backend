package api4auth

import (
	"fmt"
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/api4debtus"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/strongo/validation"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func httpSignInFromTelegramMiniapp(w http.ResponseWriter, r *http.Request) {
	log.Println("httpSignInFromTelegramMiniapp()")
	if r.Method != http.MethodPost {
		apicore.ReturnError(r.Context(), w, r, validation.NewBadRequestError(fmt.Errorf("%s method is not allowed", r.Method)))
		return
	}

	ctx := r.Context()

	// Read request body to a string
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN_ALEXTDEVBOT")

	initDataStr := string(bodyBytes)
	if err = twainitdata.Validate(initDataStr, botToken, 10*time.Second); err != nil {
		err = validation.NewBadRequestError(err)
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	var initData twainitdata.InitData
	if initData, err = twainitdata.Parse(initDataStr); err != nil {
		err = validation.NewBadRequestError(err)
		apicore.ReturnError(ctx, w, r, err)
	}

	var (
		tgBotUser facade4auth.TelegramPlatformUserEntry
		isNewUser bool
	)
	remoteClientInfo := dbmodels.RemoteClientInfo{
		HostOrApp:  r.Host,
		RemoteAddr: r.RemoteAddr,
	}
	if tgBotUser, isNewUser, err = facade4auth.SignInWithTelegram(ctx, initData, remoteClientInfo); err != nil {
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	api4debtus.ReturnToken(ctx, w, tgBotUser.Data.GetAppUserID(), telegram.PlatformID, isNewUser, false)
}
