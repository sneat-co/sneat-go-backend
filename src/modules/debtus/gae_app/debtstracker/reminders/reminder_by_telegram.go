package reminders

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/core/queues"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/platforms/tgbots"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/bot/profiles/debtus/dtb_common"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/analytics"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"time"
)

func sendReminderByTelegram(c context.Context, transfer models4debtus.TransferEntry, reminder models4debtus.Reminder, tgChatID int64, tgBot string) (sent, channelDisabledByUser bool, err error) {
	logus.Debugf(c, "sendReminderByTelegram(transfer.ContactID=%v, reminder.ContactID=%v, tgChatID=%v, tgBot=%v)", transfer.ID, reminder.ID, tgChatID, tgBot)

	if tgChatID == 0 {
		panic("tgChatID == 0")
	}
	if tgBot == "" {
		panic("tgBot is empty string")
	}

	var locale i18n.Locale

	if locale, err = facade4debtus.GetLocale(c, tgBot, tgChatID, reminder.Data.UserID); err != nil {
		return
	}

	//if !tgChat.DtForbidden.IsZero() {
	//	logus.Infof(c, "Telegram chat(id=%v) is not available since: %v", tgChatID, tgChat.DtForbidden)
	//	return false
	//}

	translator := i18n.NewSingleMapTranslator(locale, i18n.NewMapTranslator(c, trans.TRANS))

	env := dtdal.HttpAppHost.GetEnvironment(c, nil)

	if botSettings, ok := tgbots.Bots(env, nil).ByCode[tgBot]; !ok {
		err = fmt.Errorf("bot settings not found (env=%v, tgBotID=%v)", env, tgBot)
		return
	} else {
		tgBotApi := tgbotapi.NewBotAPIWithClient(botSettings.Token, dtdal.HttpClient(c))
		messageText := fmt.Sprintf(
			"<b>%v</b>\n%v\n\n",
			translator.Translate(trans.MESSAGE_TEXT_REMINDER),
			translator.Translate(trans.MESSAGE_TEXT_REMINDER_ASK_IF_RETURNED),
		)

		utm := common4debtus.UtmParams{
			Source:   "TODO",
			Medium:   telegram.PlatformID,
			Campaign: common4debtus.UTM_CAMPAIGN_REMINDER,
		}
		messageText += common4debtus.TextReceiptForTransfer(c, translator, transfer, reminder.Data.UserID, common4debtus.ShowReceiptToAutodetect, utm)

		messageConfig := tgbotapi.NewMessage(tgChatID, messageText)

		err = facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) (err error) {
			reminder, err = dtdal.Reminder.GetReminderByID(c, tx, reminder.ID)
			if err != nil {
				return err
			}
			callbackData := fmt.Sprintf(dtb_common.DEBT_RETURN_CALLBACK_DATA, dtb_common.CALLBACK_DEBT_RETURNED_PATH, reminder.ID, "%v")
			messageConfig.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					{Text: translator.Translate(trans.COMMAND_TEXT_REMINDER_RETURNED_IN_FULL), CallbackData: fmt.Sprintf(callbackData, dtb_common.RETURNED_FULLY)},
				},
				[]tgbotapi.InlineKeyboardButton{
					{Text: translator.Translate(trans.COMMAND_TEXT_REMINDER_RETURNED_PARTIALLY), CallbackData: fmt.Sprintf(callbackData, dtb_common.RETURNED_PARTIALLY)},
				},
				[]tgbotapi.InlineKeyboardButton{
					{Text: translator.Translate(trans.COMMAND_TEXT_REMINDER_NOT_RETURNED), CallbackData: fmt.Sprintf(callbackData, dtb_common.RETURNED_NOTHING)},
				},
			)
			messageConfig.ParseMode = "HTML"
			message, err := tgBotApi.Send(messageConfig)
			if err != nil {
				if _, isForbidden := err.(tgbotapi.ErrAPIForbidden); isForbidden { // TODO: Mark chat as deleted?
					logus.Infof(c, "Telegram bot API returned status 'forbidden' - either issue with token or chat deleted by user")
					if err2 := DelaySetChatIsForbidden(c, botSettings.Code, tgChatID, time.Now()); err2 != nil {
						logus.Errorf(c, "Failed to delay to set chat as forbidden: %v", err2)
					}
					channelDisabledByUser = true
					return nil // Do not pass error up
				} else {
					logus.Debugf(c, "messageConfig.Text: %v", messageConfig.Text)
					return fmt.Errorf("Failed in call to Telegram API: %w", err)
				}
			}
			sent = true
			logus.Infof(c, "Sent message to telegram. MessageID: %v", message.MessageID)

			if err = dtdal.Reminder.SetReminderIsSentInTransaction(tc, tx, reminder, time.Now(), int64(message.MessageID), "", locale.Code5, ""); err != nil {
				err = dtdal.Reminder.DelaySetReminderIsSent(tc, reminder.ID, time.Now(), int64(message.MessageID), "", locale.Code5, "")
			}
			//
			return
		}, nil)

		if err != nil {
			logus.Errorf(c, fmt.Errorf("error while sending by Telegram: %w", err).Error())
			return
		}
		if sent {
			analytics.ReminderSent(c, reminder.Data.UserID, translator.Locale().Code5, telegram.PlatformID)
		}
	}
	return
}

func DelaySetChatIsForbidden(c context.Context, botID string, tgChatID int64, at time.Time) error {
	return delaySetChatIsForbidden.EnqueueWork(c, delaying.With(queues.QueueChats, "set-chat-is-forbidden", 0), botID, tgChatID, at)
}

func SetChatIsForbidden(c context.Context, botID string, tgChatID int64, at time.Time) error {
	logus.Debugf(c, "SetChatIsForbidden(tgChatID=%v, at=%v)", tgChatID, at)
	panic("TODO: Implement SetChatIsForbidden")
	//err := gaehost.MarkTelegramChatAsForbidden(c, botID, tgChatID, at)
	//if err == nil {
	//	logus.Infof(c, "Success")
	//} else {
	//	logus.Errorf(c, err.Error())
	//	if err == datastore.ErrNoSuchEntity {
	//		return nil // Do not re-try
	//	}
	//}
	//return err
}
