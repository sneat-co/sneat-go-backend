package gaedal

import (
	"context"
	"errors"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw-telegram-models/botsfwtgmodels"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/platforms/tgbots"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade2debtus"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/general"
	"github.com/strongo/delaying"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp/appuser"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (UserDalGae) DelaySetUserPreferredLocale(c context.Context, delay time.Duration, userID string, localeCode5 string) error {
	return delaySetUserPreferredLocale.EnqueueWork(c, delaying.With(common.QUEUE_USERS, "set-user-preferred-locale", delay), userID, localeCode5)
}

func delayedSetUserPreferredLocale(c context.Context, userID string, localeCode5 string) (err error) {
	logus.Debugf(c, "delayedSetUserPreferredLocale(userID=%v, localeCode5=%v)", userID, localeCode5)
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}
	return db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		user, err := facade2debtus.User.GetUserByID(tc, tx, userID)
		if dal.IsNotFound(err) {
			logus.Errorf(c, "User not found by ID: %v", err)
			return nil
		}
		if err == nil && user.Data.PreferredLanguage != localeCode5 {
			user.Data.PreferredLanguage = localeCode5

			if err = facade2debtus.User.SaveUser(tc, tx, user); err != nil {
				err = fmt.Errorf("failed to save user to db: %w", err)
			}
		}
		return err
	}, nil)
}

func (TransferDalGae) DelayUpdateTransferWithCreatorReceiptTgMessageID(c context.Context, botCode string, transferID string, creatorTgChatID, creatorTgReceiptMessageID int64) error {
	// logus.Debugf(c, "delayUpdateTransferWithCreatorReceiptTgMessageID(botCode=%v, transferID=%v, creatorTgChatID=%v, creatorTgReceiptMessageID=%v)", botCode, transferID, creatorTgChatID, creatorTgReceiptMessageID)

	if err := delayUpdateTransferWithCreatorReceiptTgMessageID.EnqueueWork(
		c, delaying.With(common.QUEUE_TRANSFERS, "update-transfer-with-creator-receipt-tg-message-id", 0),
		botCode, transferID, creatorTgChatID, creatorTgReceiptMessageID); err != nil {
		return fmt.Errorf("failed to create delayed task update-transfer-with-creator-receipt-tg-message-id: %w", err)
	}
	return nil
}

func delayedUpdateTransferWithCreatorReceiptTgMessageID(c context.Context, botCode string, transferID string, creatorTgChatID, creatorTgReceiptMessageID int64) (err error) {
	logus.Infof(c, "delayedUpdateTransferWithCreatorReceiptTgMessageID(botCode=%v, transferID=%v, creatorTgChatID=%v, creatorReceiptTgMessageID=%v)", botCode, transferID, creatorTgChatID, creatorTgReceiptMessageID)
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}
	return db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		transfer, err := facade2debtus.Transfers.GetTransferByID(c, tx, transferID)
		if err != nil {
			logus.Errorf(c, "Failed to get transfer by ID: %v", err)
			if dal.IsNotFound(err) {
				return nil
			} else {
				return err
			}
		}
		logus.Debugf(c, "Loaded transfer: %v", transfer.Data)
		if transfer.Data.Creator().TgBotID != botCode || transfer.Data.Creator().TgChatID != creatorTgChatID || transfer.Data.CreatorTgReceiptByTgMsgID != creatorTgReceiptMessageID {
			transfer.Data.Creator().TgBotID = botCode
			transfer.Data.Creator().TgChatID = creatorTgChatID
			transfer.Data.CreatorTgReceiptByTgMsgID = creatorTgReceiptMessageID
			if err = facade2debtus.Transfers.SaveTransfer(c, tx, transfer); err != nil {
				err = fmt.Errorf("failed to save transfer to db: %w", err)
			}
		}
		return err
	}, nil)
}

func (ReceiptDalGae) DelayCreateAndSendReceiptToCounterpartyByTelegram(c context.Context, env string, transferID string, userID string) error {
	logus.Debugf(c, "delaySendReceiptToCounterpartyByTelegram(env=%v, transferID=%v, userID=%v)", env, transferID, userID)
	return delayCreateAndSendReceiptToCounterpartyByTelegram.EnqueueWork(c, delaying.With(common.QUEUE_RECEIPTS, "create-and-send-receipt-for-counterparty-by-telegram", 0), env, transferID, userID)
}

func GetTelegramChatByUserID(c context.Context, userID string) (entityID string, chat botsfwtgmodels.TgChatData, err error) {
	tgChatQuery := dal.From(botsfwtgmodels.TgChatCollection).
		WhereField("AppUserIntID", dal.Equal, userID).
		OrderBy(dal.DescendingField("DtUpdated")).
		Limit(1).
		SelectInto(models.NewDebtusTelegramChatRecord)

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}

	var tgChatRecords []dal.Record
	if tgChatRecords, err = db.QueryAllRecords(c, tgChatQuery); err != nil {
		err = fmt.Errorf("failed to load telegram chat by app user id=%v: %w", userID, err)
		return
	}
	switch len(tgChatRecords) {
	case tgChatQuery.Limit():
		entityID = fmt.Sprintf("%v", tgChatRecords[0].Key().ID)
		tgChatBase := tgChatRecords[0].Data().(models.DebtusTelegramChatData).TgChatBaseData
		chat = &tgChatBase
		return
	case 0:
		err = fmt.Errorf("%w: telegram chat not found by userID=%s:%T", dal.ErrRecordNotFound, userID, userID)
		return
	default:
		err = fmt.Errorf("%w: too many telegram chats found by userID=%s:%T: %d", dal.ErrRecordNotFound, userID, userID, len(tgChatRecords))
		return
	}
}

func DelayOnReceiptSentSuccess(c context.Context, sentAt time.Time, receiptID, transferID string, tgChatID int64, tgMsgID int, tgBotID, locale string) error {
	if receiptID == "" {
		return errors.New("receiptID == 0")
	}
	if transferID == "" {
		return errors.New("transferID == 0")
	}
	if err := delayedOnReceiptSentSuccess.EnqueueWork(c, delaying.With(common.QUEUE_RECEIPTS, "on-receipt-sent-success", 0), sentAt, receiptID, transferID, tgChatID, tgMsgID, tgBotID, locale); err != nil {
		logus.Errorf(c, err.Error())
		return onReceiptSentSuccess(c, sentAt, receiptID, transferID, tgChatID, tgMsgID, tgBotID, locale)
	}
	return nil
}

func DelayOnReceiptSendFail(c context.Context, receiptID string, tgChatID int64, tgMsgID int, failedAt time.Time, locale, details string) error {
	if receiptID == "" {
		return errors.New("receiptID == 0")
	}
	if failedAt.IsZero() {
		return errors.New("failedAt.IsZero()")
	}
	if err := delayedOnReceiptSendFail.EnqueueWork(c, delaying.With(common.QUEUE_RECEIPTS, "on-receipt-send-fail", 0), receiptID, tgChatID, tgMsgID, failedAt, locale, details); err != nil {
		logus.Errorf(c, err.Error())
		return onReceiptSendFail(c, receiptID, tgChatID, tgMsgID, failedAt, locale, details)
	}
	return nil
}

func onReceiptSentSuccess(c context.Context, sentAt time.Time, receiptID, transferID string, tgChatID int64, tgMsgID int, tgBotID, locale string) (err error) {
	logus.Debugf(c, "onReceiptSentSuccess(sentAt=%v, receiptID=%v, transferID=%v, tgChatID=%v, tgMsgID=%v tgBotID=%v, locale=%v)", sentAt, receiptID, transferID, tgChatID, tgMsgID, tgBotID, locale)
	if receiptID == "" {
		logus.Errorf(c, "receiptID == 0")
		return

	}
	if transferID == "" {
		logus.Errorf(c, "transferID == 0")
		return
	}
	if tgChatID == 0 {
		logus.Errorf(c, "tgChatID == 0")
		return
	}
	if tgMsgID == 0 {
		logus.Errorf(c, "tgMsgID == 0")
		return
	}
	var mt string
	var receipt models.ReceiptData
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		receipt := models.NewReceipt(receiptID, nil)
		transfer := models.NewTransfer(transferID, nil)
		var (
			transferEntity models.TransferData
		)
		// TODO: Replace with DAL call?
		if err := tx.GetMulti(c, []dal.Record{receipt.Record, transfer.Record}); err != nil {
			return err
		}
		if receipt.Data.TransferID != transferID {
			return errors.New("receipt.TransferID != transferID")
		}
		if receipt.Data.Status == models.ReceiptStatusSent {
			return nil
		}

		transferEntity.Counterparty().TgBotID = tgBotID
		transferEntity.Counterparty().TgChatID = tgChatID
		receipt.Data.DtSent = sentAt
		receipt.Data.Status = models.ReceiptStatusSent
		if err = tx.SetMulti(c, []dal.Record{transfer.Record, receipt.Record}); err != nil {
			return fmt.Errorf("failed to save transfer & receipt to datastore: %w", err)
		}

		if transferEntity.DtDueOn.After(time.Now()) {
			if err := dtdal.Reminder.DelayCreateReminderForTransferUser(c, transferID, transferEntity.Counterparty().UserID); err != nil {
				return fmt.Errorf("failed to delay creation of reminder for transfer counterparty: %w", err)
			}
		}
		return nil
	}); err != nil {
		mt = err.Error()
	} else {
		var translator i18n.SingleLocaleTranslator
		if translator, err = getTranslator(c, locale); err != nil {
			return
		}
		mt = translator.Translate(trans.MESSAGE_TEXT_RECEIPT_SENT_THROW_TELEGRAM)
	}

	if err = editTgMessageText(c, tgBotID, tgChatID, tgMsgID, mt); err != nil {
		errMessage := err.Error()
		err = fmt.Errorf("failed to update Telegram message (botID=%v, chatID=%v, msgID=%v): %w", tgBotID, tgChatID, tgMsgID, err)
		if strings.Contains(errMessage, "Bad Request") && strings.Contains(errMessage, " not found") {
			logMessage := logus.Errorf
			switch {
			case receipt.DtCreated.Before(time.Now().Add(-time.Hour * 24)):
				logMessage = logus.Debugf
			case receipt.DtCreated.Before(time.Now().Add(-time.Hour)):
				logMessage = logus.Infof
			case receipt.DtCreated.Before(time.Now().Add(-time.Minute)):
				logMessage = logus.Warningf
			}
			logMessage(c, err.Error())
			err = nil
		}
		return
	}
	return
}

func onReceiptSendFail(c context.Context, receiptID string, tgChatID int64, tgMsgID int, failedAt time.Time, locale, details string) (err error) {
	logus.Debugf(c, "onReceiptSendFail(receiptID=%v, failedAt=%v)", receiptID, failedAt)
	if receiptID == "" {
		return errors.New("receiptID == 0")
	}
	var receipt models.Receipt
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if receipt, err = dtdal.Receipt.GetReceiptByID(c, tx, receiptID); err != nil {
			return err
		} else if receipt.Data.DtFailed.IsZero() {
			receipt.Data.DtFailed = failedAt
			receipt.Data.Error = details
			if ndsErr := dtdal.Receipt.UpdateReceipt(c, tx, receipt); ndsErr != nil {
				logus.Errorf(c, "Failed to update Receipt with error information: %v", ndsErr) // Discard error
			}
			return err
		}
		return nil
	}, nil); err != nil {
		return
	}

	if err = editTgMessageText(c, receipt.Data.CreatedOnID, tgChatID, tgMsgID, emoji.ERROR_ICON+" Failed to send receipt: "+details); err != nil {
		logus.Errorf(c, err.Error())
		err = nil
	}
	return
}

// func getTranslatorAndTgChatID(c context.Context, userID int64) (translator i18n.SingleLocaleTranslator, tgChatID int64, err error) {
// 	var (
// 		//transfer models.TransferEntry
// 		user models.AppUser
// 	)
// 	if user, err = facade2debtus.User.GetUserByID(c, userID); err != nil {
// 		return
// 	}
// 	if user.TelegramUserID == 0 {
// 		err = errors.New("user.TelegramUserID == 0")
// 		return
// 	}
// 	var tgChat models.DebtusTelegramChat
// 	if tgChat, err = dtdal.TgChat.GetTgChatByID(c, user.TelegramUserID); err != nil {
// 		return
// 	}
// 	localeCode := tgChat.PreferredLanguage
// 	if localeCode == "" {
// 		localeCode = user.GetPreferredLocale()
// 	}
// 	if translator, err = getTranslator(c, localeCode); err != nil {
// 		return
// 	}
// 	return
// }

func getTranslator(c context.Context, localeCode string) (translator i18n.SingleLocaleTranslator, err error) {
	logus.Debugf(c, "getTranslator(localeCode=%v)", localeCode)
	return nil, errors.New("not implemented")
	//var locale i18n.Locale
	//if locale, err = shared.TheAppContext.SupportedLocales().GetLocaleByCode5(localeCode); errors.Is(err, trans.ErrUnsupportedLocale) {
	//	if locale, err = shared.TheAppContext.SupportedLocales().GetLocaleByCode5(i18n.LocaleCodeEnUS); err != nil {
	//		return
	//	}
	//}
	//translator = i18n.NewSingleMapTranslator(locale, shared.TheAppContext.GetTranslator(c))
	//return
}

func editTgMessageText(c context.Context, tgBotID string, tgChatID int64, tgMsgID int, text string) (err error) {
	msg := tgbotapi.NewEditMessageText(tgChatID, tgMsgID, "", text)
	telegramBots := tgbots.Bots(dtdal.HttpAppHost.GetEnvironment(c, nil), nil)
	botSettings, ok := telegramBots.ByCode[tgBotID]
	if !ok {
		return fmt.Errorf("Bot settings not found by tgChat.BotID=%v, out of %v items", tgBotID, len(telegramBots.ByCode))
	}
	if err = sendToTelegram(c, msg, *botSettings); err != nil {
		return
	}
	return
}

func sendToTelegram(c context.Context, msg tgbotapi.Chattable, botSettings botsfw.BotSettings) (err error) { // TODO: Merge with same in API package
	tgApi := tgbotapi.NewBotAPIWithClient(botSettings.Token, dtdal.HttpClient(c))
	if _, err = tgApi.Send(msg); err != nil {
		return
	}
	return
}

var errReceiptStatusIsNotCreated = errors.New("receipt is not in 'created' status")

func delaySendReceiptToCounterpartyByTelegram(c context.Context, receiptID int, tgChatID int64, localeCode string) error {
	return delayedSendReceiptToCounterpartyByTelegram.EnqueueWork(c, delaying.With(common.QUEUE_RECEIPTS, "send-receipt-to-counterparty-by-telegram", time.Second/10), receiptID, tgChatID, localeCode)
}

func updateReceiptStatus(c context.Context, tx dal.ReadwriteTransaction, receiptID string, expectedCurrentStatus, newStatus string) (receipt models.Receipt, err error) {

	if err = func() (err error) {
		if receipt, err = dtdal.Receipt.GetReceiptByID(c, tx, receiptID); err != nil {
			return
		}
		if receipt.Data.Status != expectedCurrentStatus {
			return errReceiptStatusIsNotCreated
		}
		receipt.Data.Status = newStatus
		if err = tx.Set(c, receipt.Record); err != nil {
			return
		}
		return
	}(); err != nil {
		err = fmt.Errorf("failed to update receipt status from %v to %v: %w", expectedCurrentStatus, newStatus, err)
	}
	return
}

func sendReceiptToCounterpartyByTelegram(c context.Context, receiptID string, tgChatID int64, localeCode string) (err error) {
	logus.Debugf(c, "delayedSendReceiptToCounterpartyByTelegram(receiptID=%v, tgChatID=%v, localeCode=%v)", receiptID, tgChatID, localeCode)

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return
	}
	if err := db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		var receipt models.Receipt

		if receipt, err = updateReceiptStatus(c, tx, receiptID, models.ReceiptStatusCreated, models.ReceiptStatusSending); err != nil {
			logus.Errorf(c, err.Error())
			err = nil // Always stop!
			return
		}

		var transfer models.TransferEntry
		if transfer, err = facade2debtus.Transfers.GetTransferByID(c, tx, receipt.Data.TransferID); err != nil {
			logus.Errorf(c, err.Error())
			if dal.IsNotFound(err) {
				err = nil
				return
			}
			return
		}

		var counterpartyUser models.AppUser

		if counterpartyUser, err = facade2debtus.User.GetUserByID(c, tx, receipt.Data.CounterpartyUserID); err != nil {
			return
		}

		var (
			tgChat         models.DebtusTelegramChat
			failedToSend   bool
			chatsForbidden bool
		)

		creatorTgChatID, creatorTgMsgID := transfer.Data.Creator().TgChatID, int(transfer.Data.CreatorTgReceiptByTgMsgID)

		var tgAccounts []appuser.AccountKey

		if tgAccounts, err = counterpartyUser.Data.GetAccounts("telegram"); err != nil {
			return err
		}
		for _, telegramAccount := range tgAccounts {
			if telegramAccount.App == "" {
				logus.Warningf(c, "User %v has account with missing bot id => %v", counterpartyUser.ID, telegramAccount.String())
				continue
			}
			var tgChatID int64
			if tgChatID, err = strconv.ParseInt(telegramAccount.ID, 10, 64); err != nil {
				logus.Errorf(c, "invalid Telegram chat ID - not an integer: %v", telegramAccount.String())
				continue
			}
			if tgChat, err = dtdal.TgChat.GetTgChatByID(c, telegramAccount.App, tgChatID); err != nil {
				logus.Errorf(c, "failed to load user's Telegram chat entity: %v", err)
				continue
			}
			if tgChat.Data.DtForbiddenLast.IsZero() {
				if err = sendReceiptToTelegramChat(c, receipt, transfer, tgChat); err != nil {
					//failedToSend = true
					if _, forbidden := err.(tgbotapi.ErrAPIForbidden); forbidden || strings.Contains(err.Error(), "Bad Request: chat not found") {
						//chatsForbidden = true
						logus.Infof(c, "Telegram chat not found or disabled (%v): %v", tgChat.ID, err)
						panic("not implemented - commented out as needs to be refactored")
						//if err2 := gaehost.MarkTelegramChatAsForbidden(c, tgChat.Data.BotID, tgChat.Data.TelegramUserID, time.Now()); err2 != nil {
						//	logus.Errorf(c, "Failed to call MarkTelegramChatAsStopped(): %v", err2.Error())
						//}
						//return nil
					}
					return
				}
				if err = DelayOnReceiptSentSuccess(c, time.Now(), receipt.ID, transfer.ID, creatorTgChatID, creatorTgMsgID, tgChat.Data.BotID, localeCode); err != nil {
					logus.Errorf(c, fmt.Errorf("failed to call DelayOnReceiptSentSuccess(): %w", err).Error())
				}
				return
			} else {
				logus.Debugf(c, "tgChat is forbidden: %v", telegramAccount.String())
			}
			break
		}

		if failedToSend { // Notify creator that receipt has not been sent
			var translator i18n.SingleLocaleTranslator
			if translator, err = getTranslator(c, localeCode); err != nil {
				return err
			}

			locale := translator.Locale()
			if chatsForbidden {
				msgTextToCreator := emoji.ERROR_ICON + translator.Translate(trans.MESSAGE_TEXT_RECEIPT_NOT_SENT_AS_COUNTERPARTY_HAS_DISABLED_TG_BOT, transfer.Data.Counterparty().ContactName)
				if err2 := DelayOnReceiptSendFail(c, receipt.ID, creatorTgChatID, creatorTgMsgID, time.Now(), translator.Locale().Code5, msgTextToCreator); err2 != nil {
					logus.Errorf(c, fmt.Errorf("failed to update receipt entity with error info: %w", err2).Error())
				}
			}
			logus.Errorf(c, "Failed to send notification to creator by Telegram (creatorTgChatID=%v, creatorTgMsgID=%v): %v", creatorTgChatID, creatorTgMsgID, err)
			msgTextToCreator := emoji.ERROR_ICON + " " + err.Error()
			if err2 := DelayOnReceiptSendFail(c, receipt.ID, creatorTgChatID, creatorTgMsgID, time.Now(), locale.Code5, msgTextToCreator); err2 != nil {
				logus.Errorf(c, fmt.Errorf("failed to update receipt entity with error info: %w", err2).Error())
			}
			err = nil
		}
		return err
	}); err != nil {
		return err
	}
	return err
}

func sendReceiptToTelegramChat(c context.Context, receipt models.Receipt, transfer models.TransferEntry, tgChat models.DebtusTelegramChat) (err error) {
	var messageToTranslate string
	switch transfer.Data.Direction() {
	case models.TransferDirectionUser2Counterparty:
		messageToTranslate = trans.TELEGRAM_RECEIPT
	case models.TransferDirectionCounterparty2User:
		messageToTranslate = trans.TELEGRAM_RECEIPT
	default:
		panic(fmt.Errorf("Unknown direction: %v", transfer.Data.Direction()))
	}

	templateData := struct {
		FromName         string
		TransferCurrency string
	}{
		FromName:         transfer.Data.Creator().ContactName,
		TransferCurrency: string(transfer.Data.Currency),
	}

	var translator i18n.SingleLocaleTranslator
	if translator, err = getTranslator(c, tgChat.Data.GetPreferredLanguage()); err != nil {
		return err
	}

	messageText, err := common.TextTemplates.RenderTemplate(c, translator, messageToTranslate, templateData)
	if err != nil {
		return err
	}
	messageText = emoji.INCOMING_ENVELOP_ICON + " " + messageText

	logus.Debugf(c, "Message: %v", messageText)

	btnViewReceiptText := emoji.CLIPBOARD_ICON + " " + translator.Translate(trans.BUTTON_TEXT_SEE_RECEIPT_DETAILS)
	btnViewReceiptData := fmt.Sprintf("view-receipt?id=%s", receipt.ID) // TODO: Pass simple digits!
	tgMessage := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: tgChat.Data.TelegramUserID,
			ReplyMarkup: tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(btnViewReceiptText, btnViewReceiptData)),
				},
			},
		},
		ParseMode:             "HTML",
		DisableWebPagePreview: true,
		Text:                  messageText,
	}

	tgBotApi := tgbots.GetTelegramBotApiByBotCode(c, tgChat.Data.BotID)

	if _, err = tgBotApi.Send(tgMessage); err != nil {
		return
	} else {
		logus.Infof(c, "Receipt %v sent to user by Telegram bot @%v", receipt.ID, tgChat.Data.BotID)
	}

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return err
	}
	err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if receipt, err = updateReceiptStatus(c, tx, receipt.ID, models.ReceiptStatusSending, models.ReceiptStatusSent); err != nil {
			logus.Errorf(c, err.Error())
			err = nil
			return
		}
		return err
	})
	return
}

func delayedCreateAndSendReceiptToCounterpartyByTelegram(c context.Context, env string, transferID string, toUserID string) error {
	logus.Debugf(c, "delayCreateAndSendReceiptToCounterpartyByTelegram(transferID=%v, toUserID=%v)", transferID, toUserID)
	if transferID == "" {
		logus.Errorf(c, "transferID == 0")
		return nil
	}
	if toUserID == "" {
		logus.Errorf(c, "toUserID == 0")
		return nil
	}
	chatEntityID, tgChat, err := GetTelegramChatByUserID(c, toUserID)
	if err != nil {
		err2 := fmt.Errorf("failed to get Telegram chat for user (id=%v): %w", toUserID, err)
		if dal.IsNotFound(err) {
			logus.Infof(c, "No telegram for user or user not found")
			return nil
		} else {
			return err2
		}
	}
	if chatEntityID == "" {
		logus.Infof(c, "No telegram for user")
		return nil
	}
	localeCode := tgChat.BaseTgChatData().PreferredLanguage
	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return err
	}
	if err = db.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		transfer, err := facade2debtus.Transfers.GetTransferByID(c, tx, transferID)
		if err != nil {
			if dal.IsNotFound(err) {
				logus.Errorf(c, err.Error())
				return nil
			}
			return fmt.Errorf("failed to get transfer by id=%v: %v", transferID, err)
		}
		if localeCode == "" {
			toUser, err := facade2debtus.User.GetUserByID(c, tx, toUserID)
			if err != nil {
				return err
			}
			localeCode = toUser.Data.GetPreferredLocale()
		}

		var translator i18n.SingleLocaleTranslator
		if translator, err = getTranslator(c, localeCode); err != nil {
			return err
		}
		locale := translator.Locale()

		var receiptID int
		receipt := models.NewReceipt("", models.NewReceiptEntity(transfer.Data.CreatorUserID, transferID, transfer.Data.Counterparty().UserID, locale.Code5, telegram.PlatformID, strconv.FormatInt(tgChat.BaseTgChatData().TelegramUserID, 10), general.CreatedOn{
			CreatedOnID:       transfer.Data.Creator().TgBotID, // TODO: Replace with method call.
			CreatedOnPlatform: transfer.Data.CreatedOnPlatform,
		}))
		if err := tx.Set(c, receipt.Record); err != nil {
			return fmt.Errorf("failed to save receipt to DB: %w", err)
		} else {
			receiptID = receipt.Record.Key().ID.(int)
		}
		if err != nil {
			return fmt.Errorf("failed to create receipt entity: %w", err)
		}
		tgChatID := tgChat.BaseTgChatData().TelegramUserID
		if err = delaySendReceiptToCounterpartyByTelegram(c, receiptID, tgChatID, localeCode); err != nil { // TODO: ideally should be called inside transaction
			logus.Errorf(c, "failed to queue receipt sending: %v", err)
			return nil
		}
		return err
	}); err != nil {
		return err
	}
	return nil
}

func (UserDalGae) DelayUpdateUserHasDueTransfers(c context.Context, userID string) error {
	if userID == "" {
		panic("userID == 0")
	}
	return delayUpdateUserHasDueTransfers.EnqueueWork(c, delaying.With(common.QUEUE_USERS, "update-user-has-due-transfers", 0), userID)
}

func delayedUpdateUserHasDueTransfers(c context.Context, userID string) (err error) {
	logus.Debugf(c, "delayUpdateUserHasDueTransfers(userID=%v)", userID)
	if userID == "" {
		logus.Errorf(c, "userID == 0")
		return nil
	}
	user, err := facade2debtus.User.GetUserByID(c, nil, userID)
	if err != nil {
		if dal.IsNotFound(err) {
			logus.Errorf(c, err.Error())
			return nil
		}
		return err
	}
	if user.Data.HasDueTransfers {
		logus.Infof(c, "Already user.HasDueTransfers == %v", user.Data.HasDueTransfers)
		return nil
	}

	q := dal.From(models.TransfersCollection).
		WhereField("BothUserIDs", dal.Equal, userID).
		WhereField("IsOutstanding", dal.Equal, true).
		WhereField("DtDueOn", dal.GreaterThen, time.Time{}).
		Limit(1).
		SelectKeysOnly(reflect.Int)

	var db dal.DB
	if db, err = GetDatabase(c); err != nil {
		return err
	}

	var reader dal.Reader
	if reader, err = db.QueryReader(c, q); err != nil {
		return err
	}

	var transferIDs []int
	transferIDs, err = dal.SelectAllIDs[int](reader, dal.WithLimit(q.Limit()))

	if len(transferIDs) > 0 {
		// panic("Not implemented - refactoring in progress")
		// reminder := reminders[0]
		err = db.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
			if user, err := facade2debtus.User.GetUserByID(tc, tx, userID); err != nil {
				if dal.IsNotFound(err) {
					logus.Errorf(c, err.Error())
					return nil // Do not retry
				}
				return err
			} else if !user.Data.HasDueTransfers {
				user.Data.HasDueTransfers = true
				if err = tx.Set(tc, user.Record); err != nil {
					return fmt.Errorf("failed to save user to db: %w", err)
				}
				logus.Infof(c, "User updated & saved to datastore")
			}
			return nil
		}, nil)
	}
	return err
}
