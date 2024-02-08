package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/crediterra/money"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/platforms/tgbots"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_transfer"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/auth"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func handleError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(err.Error()))
}
func handleTgHelperCurrencySelected(c context.Context, w http.ResponseWriter, r *http.Request, authInfo auth.AuthInfo) {
	if err := r.ParseForm(); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
	selectedCurrency := r.FormValue("currency")
	if selectedCurrency == "" {
		handleError(w, http.StatusBadRequest, errors.New("missing required parameter 'currency'"))
		return
	}
	if len(selectedCurrency) != 3 {
		handleError(w, http.StatusBadRequest, errors.New("wrong lengths of parameter 'currency'"))
		return
	}
	if strings.ToUpper(selectedCurrency) != selectedCurrency {
		handleError(w, http.StatusBadRequest, errors.New("wrong currency code"))
		return
	}

	tgChatKeyID := r.Form.Get("tg-chat")
	if tgChatKeyID == "" {
		handleError(w, http.StatusBadRequest, errors.New("missing required parameter chat ID"))
		return
	}

	if !strings.Contains(tgChatKeyID, ":") {
		handleError(w, http.StatusBadRequest, errors.New("wrong format of Telegram chat ID parameter"))
		return
	}

	tgChatID, err := strconv.ParseInt(strings.Split(tgChatKeyID, ":")[1], 10, 64)
	if err != nil {
		handleError(w, http.StatusBadRequest, fmt.Errorf("value of Telegram chat ID should be integer: %w", err))
		return
	}
	log.Debugf(c, "AppUserIntID: %v, tgChatKeyID: %v", authInfo.UserID, tgChatKeyID)

	errs := make(chan error, 2) // We use errors channel as sync pipe

	var user models.AppUser

	var userTask sync.WaitGroup

	userTask.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf(c, "panic in handleTgHelperCurrencySelected() => dtdal.User.SetLastCurrency(): %v", r)
			}
		}()
		if err := dtdal.User.SetLastCurrency(c, authInfo.UserID, money.CurrencyCode(selectedCurrency)); err != nil {
			log.Errorf(c, "Failed to save user last currency: %v", err)
		}
		userTask.Done()
		errs <- nil
	}()

	go func(currency string) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf(c, "panic in handleTgHelperCurrencySelected() => dtdal.TgChat.DoSomething() => sendToTelegram(): %v", r)
			}
		}()
		errs <- dtdal.TgChat.DoSomething(c, &userTask, currency, tgChatID, authInfo, user,
			func(tgChat botsfwtgmodels.TgChatData) error {
				// TODO: This is some serious architecture sheet. Too sleepy to make it right, just make it working.
				return sendToTelegram(c, user, tgChatID, tgChat, &userTask, r)
			},
		)
	}(selectedCurrency)

	for i := range []int{1, 2} {
		if err := <-errs; err != nil {
			log.Errorf(c, "%v: %v", i, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}

	log.Debugf(c, "Selected currency processed: %v", selectedCurrency)
}

// TODO: This is some serious architecture sheet. Too sleepy to make it right, just make it working.
func sendToTelegram(c context.Context, user models.AppUser, tgChatID int64, tgChat botsfwtgmodels.TgChatData, userTask *sync.WaitGroup, r *http.Request) (err error) {
	telegramBots := tgbots.Bots(dtdal.HttpAppHost.GetEnvironment(c, nil), nil)
	baseChatData := tgChat.BaseTgChatData()
	botID := baseChatData.BotID
	botSettings, ok := telegramBots.ByCode[botID]
	if !ok {
		return fmt.Errorf("ReferredTo settings not found by tgChat.BotID=%v, out of %v items", botID, len(telegramBots.ByCode))
	}

	log.Debugf(c, "botSettings(%v : %v)", botSettings.Code, botSettings.Token)

	tgBotApi := tgbotapi.NewBotAPIWithClient(botSettings.Token, dtdal.HttpClient(c))
	tgBotApi.EnableDebug(c)

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
