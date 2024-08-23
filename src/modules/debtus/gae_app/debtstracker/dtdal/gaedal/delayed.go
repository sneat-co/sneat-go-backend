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
	"github.com/sneat-co/sneat-go-backend/src/auth/facade4auth"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/const4debtus"
	tgbots2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/platforms/debtustgbots"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/delaying"
	"github.com/strongo/i18n"
	"github.com/strongo/logus"
	"github.com/strongo/strongoapp/appuser"
	"strconv"
	"strings"
	"time"
)

func (TransferDalGae) DelayUpdateTransferWithCreatorReceiptTgMessageID(c context.Context, botCode string, transferID string, creatorTgChatID, creatorTgReceiptMessageID int64) error {
	// logus.Debugf(c, "delayerUpdateTransferWithCreatorReceiptTgMessageID(botCode=%v, transferID=%v, creatorTgChatID=%v, creatorTgReceiptMessageID=%v)", botCode, transferID, creatorTgChatID, creatorTgReceiptMessageID)

	if err := delayerUpdateTransferWithCreatorReceiptTgMessageID.EnqueueWork(
		c, delaying.With(const4debtus.QueueTransfers, "update-transfer-with-creator-receipt-tg-message-id", 0),
		botCode, transferID, creatorTgChatID, creatorTgReceiptMessageID); err != nil {
		return fmt.Errorf("failed to create delayed task update-transfer-with-creator-receipt-tg-message-id: %w", err)
	}
	return nil
}

func delayedUpdateTransferWithCreatorReceiptTgMessageID(c context.Context, botCode string, transferID string, creatorTgChatID, creatorTgReceiptMessageID int64) (err error) {
	logus.Infof(c, "delayedUpdateTransferWithCreatorReceiptTgMessageID(botCode=%v, transferID=%v, creatorTgChatID=%v, creatorReceiptTgMessageID=%v)", botCode, transferID, creatorTgChatID, creatorTgReceiptMessageID)
	return facade.RunReadwriteTransaction(c, func(tc context.Context, tx dal.ReadwriteTransaction) error {
		transfer, err := facade4debtus.Transfers.GetTransferByID(c, tx, transferID)
		if err != nil {
			logus.Errorf(c, "Failed to get transfer by ContactID: %v", err)
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
			if err = facade4debtus.Transfers.SaveTransfer(c, tx, transfer); err != nil {
				err = fmt.Errorf("failed to save transfer to db: %w", err)
			}
		}
		return err
	}, nil)
}

func (ReceiptDalGae) DelayCreateAndSendReceiptToCounterpartyByTelegram(c context.Context, env string, transferID string, userID string) error {
	logus.Debugf(c, "delayerSendReceiptToCounterpartyByTelegram(env=%v, transferID=%v, userID=%v)", env, transferID, userID)
	return delayerCreateAndSendReceiptToCounterpartyByTelegram.EnqueueWork(c, delaying.With(const4debtus.QueueReceipts, "create-and-send-receipt-for-counterparty-by-telegram", 0), env, transferID, userID)
}

func GetTelegramChatByUserID(c context.Context, userID string) (entityID string, chat botsfwtgmodels.TgChatData, err error) {
	tgChatQuery := dal.From(botsfwtgmodels.TgChatCollection).
		WhereField("AppUserIntID", dal.Equal, userID).
		OrderBy(dal.DescendingField("DtUpdated")).
		Limit(1).
		SelectInto(models4debtus.NewDebtusTelegramChatRecord)

	var db dal.DB
	if db, err = facade.GetDatabase(c); err != nil {
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
		tgChatBase := tgChatRecords[0].Data().(models4debtus.DebtusTelegramChatData).TgChatBaseData
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
	if err := delayerOnReceiptSentSuccess.EnqueueWork(c, delaying.With(const4debtus.QueueReceipts, "on-receipt-sent-success", 0), sentAt, receiptID, transferID, tgChatID, tgMsgID, tgBotID, locale); err != nil {
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
	if err := delayerOnReceiptSendFail.EnqueueWork(c, delaying.With(const4debtus.QueueReceipts, "on-receipt-send-fail", 0), receiptID, tgChatID, tgMsgID, failedAt, locale, details); err != nil {
		logus.Errorf(c, err.Error())
		return delayedOnReceiptSendFail(c, receiptID, tgChatID, tgMsgID, failedAt, locale, details)
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
	var receipt models4debtus.ReceiptDbo
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		receipt := models4debtus.NewReceipt(receiptID, nil)
		transfer := models4debtus.NewTransfer(transferID, nil)
		var (
			transferEntity models4debtus.TransferData
		)
		// TODO: Replace with DAL call?
		if err := tx.GetMulti(c, []dal.Record{receipt.Record, transfer.Record}); err != nil {
			return err
		}
		if receipt.Data.TransferID != transferID {
			return errors.New("receipt.TransferID != transferID")
		}
		if receipt.Data.Status == models4debtus.ReceiptStatusSent {
			return nil
		}

		transferEntity.Counterparty().TgBotID = tgBotID
		transferEntity.Counterparty().TgChatID = tgChatID
		receipt.Data.DtSent = sentAt
		receipt.Data.Status = models4debtus.ReceiptStatusSent
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

func delayedOnReceiptSendFail(c context.Context, receiptID string, tgChatID int64, tgMsgID int, failedAt time.Time, locale, details string) (err error) {
	logus.Debugf(c, "delayedOnReceiptSendFail(receiptID=%v, failedAt=%v)", receiptID, failedAt)
	if receiptID == "" {
		return errors.New("receiptID == 0")
	}
	var receipt models4debtus.ReceiptEntry
	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		if receipt, err = dtdal.Receipt.GetReceiptByID(c, tx, receiptID); err != nil {
			return err
		} else if receipt.Data.DtFailed.IsZero() {
			receipt.Data.DtFailed = failedAt
			receipt.Data.Error = details
			if ndsErr := dtdal.Receipt.UpdateReceipt(c, tx, receipt); ndsErr != nil {
				logus.Errorf(c, "Failed to update ReceiptEntry with error information: %v", ndsErr) // Discard error
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
// 		user models.AppUserOBSOLETE
// 	)
// 	if user, err = dal4userus.GetUserByID(c, userID); err != nil {
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
	//if locale, err = anybot.TheAppContext.SupportedLocales().GetLocaleByCode5(localeCode); errors.Is(err, trans.ErrUnsupportedLocale) {
	//	if locale, err = anybot.TheAppContext.SupportedLocales().GetLocaleByCode5(i18n.LocaleCodeEnUS); err != nil {
	//		return
	//	}
	//}
	//translator = i18n.NewSingleMapTranslator(locale, anybot.TheAppContext.GetTranslator(c))
	//return
}

func editTgMessageText(c context.Context, tgBotID string, tgChatID int64, tgMsgID int, text string) (err error) {
	msg := tgbotapi.NewEditMessageText(tgChatID, tgMsgID, "", text)
	telegramBots := tgbots2.Bots(dtdal.HttpAppHost.GetEnvironment(c, nil))
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

func DelaySendReceiptToCounterpartyByTelegram(c context.Context, receiptID string, tgChatID int64, localeCode string) error {
	return delayerSendReceiptToCounterpartyByTelegram.EnqueueWork(c, delaying.With(const4debtus.QueueReceipts, "send-receipt-to-counterparty-by-telegram", time.Second/10), receiptID, tgChatID, localeCode)
}

func updateReceiptStatus(c context.Context, tx dal.ReadwriteTransaction, receiptID string, expectedCurrentStatus, newStatus string) (receipt models4debtus.ReceiptEntry, err error) {

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

func delayedSendReceiptToCounterpartyByTelegram(c context.Context, receiptID string, tgChatID int64, localeCode string) (err error) {
	logus.Debugf(c, "delayedSendReceiptToCounterpartyByTelegram(receiptID=%v, tgChatID=%v, localeCode=%v)", receiptID, tgChatID, localeCode)

	if err := facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		var receipt models4debtus.ReceiptEntry

		if receipt, err = updateReceiptStatus(c, tx, receiptID, models4debtus.ReceiptStatusCreated, models4debtus.ReceiptStatusSending); err != nil {
			logus.Errorf(c, err.Error())
			err = nil // Always stop!
			return
		}

		var transfer models4debtus.TransferEntry
		if transfer, err = facade4debtus.Transfers.GetTransferByID(c, tx, receipt.Data.TransferID); err != nil {
			logus.Errorf(c, err.Error())
			if dal.IsNotFound(err) {
				err = nil
				return
			}
			return
		}

		counterpartyUser := dbo4userus.NewUserEntry(receipt.Data.CounterpartyUserID)
		if err = dal4userus.GetUser(c, tx, counterpartyUser); err != nil {
			return
		}

		var (
			tgChat         models4debtus.DebtusTelegramChat
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
				logus.Warningf(c, "UserEntry %v has account with missing bot id => %v", counterpartyUser.ID, telegramAccount.String())
				continue
			}
			var tgChatID int64
			if tgChatID, err = strconv.ParseInt(telegramAccount.ID, 10, 64); err != nil {
				logus.Errorf(c, "invalid Telegram chat ContactID - not an integer: %v", telegramAccount.String())
				continue
			}
			if tgChat, err = facade4auth.TgChat.GetTgChatByID(c, telegramAccount.App, tgChatID); err != nil {
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
				if err = DelayOnReceiptSentSuccess(c, time.Now(), receipt.ID, transfer.ID, creatorTgChatID, creatorTgMsgID, tgChat.Key.Parent().ID.(string), localeCode); err != nil {
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

func sendReceiptToTelegramChat(c context.Context, receipt models4debtus.ReceiptEntry, transfer models4debtus.TransferEntry, tgChat models4debtus.DebtusTelegramChat) (err error) {
	var messageToTranslate string
	switch transfer.Data.Direction() {
	case models4debtus.TransferDirectionUser2Counterparty:
		messageToTranslate = trans.TELEGRAM_RECEIPT
	case models4debtus.TransferDirectionCounterparty2User:
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

	messageText, err := common4debtus.TextTemplates.RenderTemplate(c, translator, messageToTranslate, templateData)
	if err != nil {
		return err
	}
	messageText = emoji.INCOMING_ENVELOP_ICON + " " + messageText

	logus.Debugf(c, "Message: %v", messageText)

	btnViewReceiptText := emoji.CLIPBOARD_ICON + " " + translator.Translate(trans.BUTTON_TEXT_SEE_RECEIPT_DETAILS)
	btnViewReceiptData := fmt.Sprintf("view-receipt?id=%s", receipt.ID) // TODO: Pass simple digits!

	var telegramUserID int64
	if telegramUserID, err = strconv.ParseInt(tgChat.Data.BotUserIDs[0], 10, 64); err != nil {
		return err
	}
	tgMessage := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: telegramUserID,
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

	tgBotApi := tgbots2.GetTelegramBotApiByBotCode(c, tgChat.Key.Parent().ID.(string))

	if _, err = tgBotApi.Send(tgMessage); err != nil {
		return
	} else {
		logus.Infof(c, "ReceiptEntry %v sent to user by Telegram bot @%v", receipt.ID, tgChat.Key.Parent().ID.(string))
	}

	err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		if receipt, err = updateReceiptStatus(c, tx, receipt.ID, models4debtus.ReceiptStatusSending, models4debtus.ReceiptStatusSent); err != nil {
			logus.Errorf(c, err.Error())
			err = nil
			return
		}
		return err
	})
	return
}

func delayedCreateAndSendReceiptToCounterpartyByTelegram(c context.Context, env string, transferID string, toUserID string) error {
	logus.Debugf(c, "delayerCreateAndSendReceiptToCounterpartyByTelegram(transferID=%v, toUserID=%v)", transferID, toUserID)
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

	if err = facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) error {
		var transfer models4debtus.TransferEntry
		transfer, err = facade4debtus.Transfers.GetTransferByID(c, tx, transferID)
		if err != nil {
			if dal.IsNotFound(err) {
				logus.Errorf(c, err.Error())
				return nil
			}
			return fmt.Errorf("failed to get transfer by id=%v: %v", transferID, err)
		}
		if localeCode == "" {
			toUser, err := dal4userus.GetUserByID(c, tx, toUserID)
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

		var receiptID string
		receipt := models4debtus.NewReceipt("", models4debtus.NewReceiptEntity(transfer.Data.CreatorUserID, transferID, transfer.Data.Counterparty().UserID, locale.Code5, telegram.PlatformID, tgChat.BaseTgChatData().BotUserIDs[0], general.CreatedOn{
			CreatedOnID:       transfer.Data.Creator().TgBotID, // TODO: Replace with method call.
			CreatedOnPlatform: transfer.Data.CreatedOnPlatform,
		}))
		if err := tx.Set(c, receipt.Record); err != nil {
			return fmt.Errorf("failed to save receipt to DB: %w", err)
		} else {
			receiptID = receipt.Record.Key().ID.(string)
		}
		if err != nil {
			return fmt.Errorf("failed to create receipt entity: %w", err)
		}
		var tgChatID int64
		if tgChatID, err = strconv.ParseInt(tgChat.BaseTgChatData().BotUserIDs[0], 10, 64); err != nil {
			return err
		}
		if err = DelaySendReceiptToCounterpartyByTelegram(c, receiptID, tgChatID, localeCode); err != nil { // TODO: ideally should be called inside transaction
			logus.Errorf(c, "failed to queue receipt sending: %v", err)
			return nil
		}
		return err
	}); err != nil {
		return err
	}
	return nil
}
