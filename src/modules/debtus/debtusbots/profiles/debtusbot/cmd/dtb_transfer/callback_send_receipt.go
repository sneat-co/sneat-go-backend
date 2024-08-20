package dtb_transfer

import (
	"context"
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/strongo/logus"
	"html"
	"net/url"
	"strings"

	"errors"
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/sneat-co/debtstracker-translations/emoji"
)

var SendReceiptCallbackCommand = botsfw.NewCallbackCommand(SendReceiptCallbackPath, CallbackSendReceipt)

func CallbackSendReceipt(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	q := callbackUrl.Query()
	sendBy := q.Get("by")
	spaceID := q.Get("spaceID")
	logus.Debugf(c, "CallbackSendReceipt(callbackUrl=%v)", callbackUrl)
	return m, facade.RunReadwriteTransaction(c, func(c context.Context, tx dal.ReadwriteTransaction) (err error) {
		var (
			transferID string
			transfer   models4debtus.TransferEntry
		)
		transferID = q.Get(WizardParamTransfer)
		if transferID == "" {
			return fmt.Errorf("missing transfer ContactID")
		}
		transfer, err = facade4debtus.Transfers.GetTransferByID(c, tx, transferID)
		if err != nil {
			return fmt.Errorf("failed to get transfer by ContactID: %w", err)
		}
		//chatEntity := whc.ChatData() //TODO: Need this to get appUser, has to be refactored
		//appUser, err := whc.GetAppUser()
		counterparty, err := facade4debtus.GetDebtusSpaceContactByID(c, tx, spaceID, transfer.Data.Counterparty().ContactID)
		if err != nil {
			return err
		}
		if IsTransferNotificationsBlockedForChannel(counterparty.Data, sendBy) {
			m = whc.NewMessage(trans.MESSAGE_TEXT_USER_BLOCKED_TRANSFER_NOTIFICATIONS_BY)
			return err
		}
		chatEntity := whc.ChatData()
		switch sendBy {
		case SendReceiptByChooseChannel:
			m, err = createSendReceiptOptionsMessage(whc, transfer)
			return
		case ReceiptActionDoNotSend:
			logus.Debugf(c, "CallbackSendReceipt(): do-not-send")
			if m, err = whc.NewEditMessage(whc.Translate(trans.MESSAGE_TEXT_RECEIPT_WILL_NOT_BE_SENT), botsfw.MessageFormatHTML); err != nil {
				return
			}

			// TODO: do type assertion with botsfw.CallbackQuery interface
			callbackMessage := whc.Input().(telegram.TgWebhookCallbackQuery).TelegramCallbackMessage()
			if callbackMessage != nil && callbackMessage.Text == m.Text {
				m.Text += " (double clicked)"
			}
			m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
				[]tgbotapi.InlineKeyboardButton{
					{
						Text:         whc.Translate(trans.COMMAND_TEXT_I_HAVE_CHANGED_MY_MIND),
						CallbackData: fmt.Sprintf("%v?by=%v&%v=%v", SendReceiptCallbackPath, SendReceiptByChooseChannel, WizardParamTransfer, transferID),
					},
				},
			)
			return err
		case string(models4debtus.InviteByTelegram):
			panic(fmt.Sprintf("Unsupported option: %v", models4debtus.InviteByTelegram))
		case string(models4debtus.InviteByLinkToTelegram):
			m, err = showLinkForReceiptInTelegram(whc, transfer)
			return err
		case string(models4debtus.InviteBySms):

			if counterparty.Data.PhoneNumber > 0 {
				m, err = sendReceiptBySms(whc, tx, spaceID, counterparty.Data.PhoneContact, transfer, counterparty)
				return err
			} else {
				var updateMessage botsfw.MessageFromBot
				if updateMessage, err = whc.NewEditMessage(whc.Translate(trans.MESSAGE_TEXT_LETS_SEND_SMS), botsfw.MessageFormatHTML); err != nil {
					return
				}
				if _, err = whc.Responder().SendMessage(c, updateMessage, botsfw.BotAPISendMessageOverHTTPS); err != nil {
					logus.Errorf(c, fmt.Errorf("failed to update Telegram message: %w", err).Error())
					err = nil
				}

				chatEntity.SetAwaitingReplyTo(ASK_PHONE_NUMBER_FOR_RECEIPT_COMMAND)
				chatEntity.AddWizardParam(WizardParamTransfer, transferID)
				mt := strings.Join([]string{
					whc.Translate(trans.MESSAGE_TEXT_ASK_PHONE_NUMBER_OF_COUNTERPARTY, html.EscapeString(transfer.Data.Counterparty().ContactName)),
					whc.Translate(trans.MESSAGE_TEXT_USE_CONTACT_TO_SEND_PHONE_NUMBER, emoji.PAPERCLIP_ICON),
					whc.Translate(trans.MESSAGE_TEXT_ABOUT_PHONE_NUMBER_FORMAT),
					whc.Translate(trans.MESSAGE_TEXT_THIS_NUMBER_WILL_BE_USED_TO_SEND_RECEIPT),
				}, "\n\n")
				//mt += "\n\n" + whc.Translate(trans.MESSAGE_TEXT_VIEW_MY_NUMBER_IN_INTERNATIONAL_FORMAT)

				m = whc.NewMessage(mt)
				m.Format = botsfw.MessageFormatHTML
				keyboard := [][]tgbotapi.KeyboardButton{
					{
						{RequestContact: true, Text: whc.Translate(trans.COMMAND_TEXT_VIEW_MY_NUMBER_IN_INTERNATIONAL_FORMAT)},
					},
				}
				lastName := whc.GetSender().GetLastName()
				if lastName == "Trakhimenok" || lastName == "Paltseva" {
					for k := range common4debtus.TwilioTestNumbers {
						keyboard = append(keyboard, []tgbotapi.KeyboardButton{{Text: k}})

					}
				}
				m.Keyboard = &tgbotapi.ReplyKeyboardMarkup{
					Keyboard: keyboard,
				}
			}
		case string(models4debtus.InviteByEmail):
			chatEntity.SetAwaitingReplyTo(ASK_EMAIL_FOR_RECEIPT_COMMAND)
			chatEntity.AddWizardParam(WizardParamTransfer, transferID)
			m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_INVITE_ASK_EMAIL_FOR_RECEIPT, transfer.Data.Counterparty().ContactName))
		default:
			err = errors.New("Unknown channel to send receipt: " + sendBy)
			logus.Errorf(c, err.Error())
		}
		return err
	})
}

func showLinkForReceiptInTelegram(whc botsfw.WebhookContext, transfer models4debtus.TransferEntry) (m botsfw.MessageFromBot, err error) {
	receiptData := models4debtus.NewReceiptEntity(whc.AppUserID(), transfer.ID, transfer.Data.Counterparty().UserID, whc.Locale().Code5, "link", "telegram", general.CreatedOn{
		CreatedOnPlatform: whc.BotPlatform().ID(),
		CreatedOnID:       whc.GetBotCode(),
	})
	var receipt models4debtus.ReceiptEntry
	if receipt, err = dtdal.Receipt.CreateReceipt(whc.Context(), receiptData); err != nil {
		return m, err
	}
	receiptUrl := GetUrlForReceiptInTelegram(whc.GetBotCode(), receipt.ID, whc.Locale().Code5)
	m.Text = "Send this link to counterparty:\n\n" + fmt.Sprintf(`<a href="%v">%v</a>`, receiptUrl, receiptUrl) + "\n\nPlease be aware that the first person opening this link will be treated as counterparty for this debt."
	m.Format = botsfw.MessageFormatHTML
	m.IsEdit = true
	return
}

func IsTransferNotificationsBlockedForChannel(counterparty *models4debtus.DebtusSpaceContactDbo, channel string) bool {
	for _, blockedBy := range counterparty.NoTransferUpdatesBy {
		if blockedBy == channel {
			return true
		}
	}
	return false
}
