package api4auth

import (
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-fw-telegram-webapp/twainitdata"
	"github.com/sneat-co/sneat-go-core/apicore"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"github.com/strongo/validation"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func httpLoginFromTelegramMiniapp(w http.ResponseWriter, r *http.Request) {
	log.Println("httpLoginFromTelegramMiniapp()")
	if r.Method != http.MethodPost {
		apicore.ReturnError(r.Context(), w, r, validation.NewBadRequestError(fmt.Errorf("%s method is not allowed", r.Method)))
		return
	}

	ctx := r.Context()

	if !httpserver.AccessControlAllowOrigin(w, r) {
		return
	}

	// Read request body to a string
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)
		apicore.ReturnError(ctx, w, r, err)
		return
	}

	log.Printf("Request body: \n%s\n", string(bodyBytes))

	botID := os.Getenv("TELEGRAM_BOT_ID")
	if botID == "" {
		botID = "SneatBot"
	}
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN_" + strings.ToUpper(botID))

	initDataStr := string(bodyBytes)
	if err = twainitdata.Validate(initDataStr, botToken, 10*time.Second); err != nil {
		if !(errors.Is(err, twainitdata.ErrExpired) && (strings.Contains(r.Host, ".ngrok.") || strings.HasPrefix(r.Host, "localhost:"))) {
			err = validation.NewBadRequestError(err)
			apicore.ReturnError(ctx, w, r, err)
			return
		}
	}

	var initData twainitdata.InitData
	if initData, err = twainitdata.Parse(initDataStr); err != nil {
		err = validation.NewBadRequestError(err)
		apicore.ReturnError(ctx, w, r, err)
	}

	signInWithTelegram(ctx, w, r, botID, initData)
}
