package unsorted

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/auth/token4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/platforms/debtustgbots"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_transfer"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/facade4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/httpserver"
	"github.com/strongo/logus"
	"github.com/strongo/validation"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func HandleTgHelperCurrencySelected(ctx context.Context, w http.ResponseWriter, r *http.Request, authInfo token4auth.AuthInfo) {
	if err := r.ParseForm(); err != nil {
		httpserver.HandleError(ctx, err, "", w, r)
		return
	}
	selectedCurrency := r.FormValue("currency")
	if selectedCurrency == "" {
		httpserver.HandleError(ctx, validation.NewErrRequestIsMissingRequiredField("currency"), "", w, r)
		return
	}
	if len(selectedCurrency) != 3 {
		httpserver.HandleError(ctx, validation.NewErrBadRequestFieldValue("currency", "wrong lengths of parameter value"), "", w, r)
		return
	}
	if strings.ToUpper(selectedCurrency) != selectedCurrency {
		httpserver.HandleError(ctx, validation.NewErrBadRequestFieldValue("currency", "wrong currency code"), "", w, r)
		return
	}

	tgChatKeyID := r.Form.Get("tg-chat")
	if tgChatKeyID == "" {
		httpserver.HandleError(ctx, validation.NewErrRequestIsMissingRequiredField("tg-chat"), "", w, r)
		return
	}

	if !strings.Contains(tgChatKeyID, ":") {
		httpserver.HandleError(ctx, validation.NewErrBadRequestFieldValue("tg-chat", "wrong format of Telegram chat ContactID parameter"), "", w, r)
		return
	}

	tgChatID, err := strconv.ParseInt(strings.Split(tgChatKeyID, ":")[1], 10, 64)
	if err != nil {
		httpserver.HandleError(ctx, validation.NewErrBadRequestFieldValue("tg-chat", "value of Telegram chat ContactID should be integer"), "", w, r)
		return
	}
	logus.Debugf(ctx, "AppUserIntID: %v, tgChatKeyID: %v", authInfo.UserID, tgChatKeyID)

	errs := make(chan error, 2) // We use errors channel as sync pipe

	var user dbo4userus.UserEntry

	var userTask sync.WaitGroup

	userTask.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logus.Errorf(ctx, "panic in HandleTgHelperCurrencySelected() => dtdal.UserEntry.SetLastCurrency(): %v", r)
			}
		}()
		if err := facade4userus.SetLastCurrency(ctx, facade.NewUserContext(authInfo.UserID), money.CurrencyCode(selectedCurrency)); err != nil {
			logus.Errorf(ctx, "Failed to save user last currency: %v", err)
		}
		userTask.Done()
		errs <- nil
	}()

	go func(currency string) {
		defer func() {
			if r := recover(); r != nil {
				logus.Errorf(ctx, "panic in HandleTgHelperCurrencySelected() => dtdal.TgChat.DoSomething() => sendToTelegram(): %v", r)
			}
		}()
		errs <- facade4auth.TgChat.DoSomething(ctx, &userTask, currency, tgChatID, authInfo, user,
			func(tgChat botsfwtgmodels.TgChatData) error {
				// TODO: This is some serious architecture sheet. Too sleepy to make it right, just make it working.
				botID := "TODO:setup_bot_id"
				return sendToTelegram(ctx, user, botID, tgChatID, tgChat, &userTask, r)
			},
		)
	}(selectedCurrency)

	for i := range []int{1, 2} {
		if err := <-errs; err != nil {
			logus.Errorf(ctx, "%v: %v", i, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}

	logus.Debugf(ctx, "Selected currency processed: %v", selectedCurrency)
}

// TODO: This is some serious architecture sheet. Too sleepy to make it right, just make it working.
func sendToTelegram(ctx context.Context, user dbo4userus.UserEntry, botID string, tgChatID int64, tgChat botsfwtgmodels.TgChatData, userTask *sync.WaitGroup, r *http.Request) (err error) {
	telegramBots := debtustgbots.Bots(dtdal.HttpAppHost.GetEnvironment(ctx, nil))
	baseChatData := tgChat.BaseTgChatData()
	botSettings, ok := telegramBots.ByCode[botID]
	if !ok {
		return fmt.Errorf("ReferredTo settings not found by tgChat.BotID=%v, out of %v items", botID, len(telegramBots.ByCode))
	}

	logus.Debugf(ctx, "botSettings(%v : %v)", botSettings.Code, botSettings.Token)

	tgBotApi := tgbotapi.NewBotAPIWithClient(botSettings.Token, dtdal.HttpClient(ctx))
	tgBotApi.EnableDebug(ctx)

	userTask.Wait()

	whc := NewApiWebhookContext(
		r,
		user.Data,
		user.ID,
		tgChatID,
		baseChatData,
	)

	var messageFromBot botsfw.MessageFromBot
	switch {
	case strings.Contains(baseChatData.AwaitingReplyTo, "lending"):
		messageFromBot, err = dtb_transfer.AskLendingAmountCommand.Action(whc)
	case strings.Contains(baseChatData.AwaitingReplyTo, "borrowing"):
		messageFromBot, err = dtb_transfer.AskBorrowingAmountCommand.Action(whc)
	default:
		return fmt.Errorf("tgChat.AwaitingReplyTo has unexpected value: %v", baseChatData.AwaitingReplyTo)
	}
	if err != nil {
		return fmt.Errorf("failed to create message from bot: %w", err)
	}

	messageConfig := tgbotapi.NewMessage(tgChatID, messageFromBot.Text)
	messageConfig.ReplyMarkup = messageFromBot.Keyboard
	messageConfig.ParseMode = "HTML"

	if _, err = tgBotApi.Send(messageConfig); err != nil {
		return fmt.Errorf("failed to send message to Telegram chat=%d: %w", tgChatID, err)
	}
	return nil
}
