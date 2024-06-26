package dtb_transfer

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/logus"
	"html"
	"net/url"

	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"strings"
)

//func CancelReceiptAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
//	return whc.NewMessage("TODO: Sorry, cancel is not implemented yet..."), nil
//}

const VIEW_RECEIPT_CALLBACK_COMMAND = "view-receipt"

var ViewReceiptCallbackCommand = botsfw.NewCallbackCommand(VIEW_RECEIPT_CALLBACK_COMMAND, viewReceiptCallbackAction)

func ShowReceipt(whc botsfw.WebhookContext, receiptID string) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()

	var receipt models.Receipt
	if receipt, err = dtdal.Receipt.GetReceiptByID(c, nil, receiptID); err != nil {
		return m, err
	}

	if receipt.Data.CreatorUserID == whc.AppUserID() {
		m.Text = whc.Translate(trans.MESSAGE_TEXT_RECEIPT_ATTEMPT_TO_VIEW_OWN)
		return
	}

	receipt, err = facade.MarkReceiptAsViewed(c, receiptID, whc.AppUserID())
	if err != nil {
		return
	}

	transfer, err := facade.Transfers.GetTransferByID(c, nil, receipt.Data.TransferID)
	if err != nil {
		return m, err
	}

	m = whc.NewMessage("")

	var (
		mt           string
		counterparty models.ContactEntry
	)
	counterpartyCounterparty := transfer.Data.Creator()

	if counterpartyCounterparty.ContactID != "" {
		counterparty, err = facade.GetContactByID(c, nil, counterpartyCounterparty.ContactID)
	} else {
		if user, err := facade.User.GetUserByID(c, nil, transfer.Data.CreatorUserID); err != nil {
			return m, err
		} else {
			counterparty.Data = &models.DebtusContactDbo{}
			counterparty.Data.FirstName = user.Data.FirstName
			counterparty.Data.LastName = user.Data.LastName
		}
	}

	if err != nil {
		return m, err
	}
	utm := common.NewUtmParams(whc, common.UTM_CAMPAIGN_REMINDER)
	mt = common.TextReceiptForTransfer(c, whc, transfer, whc.AppUserID(), common.ShowReceiptToAutodetect, utm)

	logus.Debugf(c, "Receipt text: %v", mt)

	var inlineKeyboard *tgbotapi.InlineKeyboardMarkup

	if receipt.Data.CreatorUserID == whc.AppUserID() {
		mt += "\n" + whc.Translate(trans.MESSAGE_TEXT_SELF_ACKNOWLEDGEMENT, html.EscapeString(transfer.Data.Counterparty().ContactName))
	} else {
		isAcknowledgedAlready := !transfer.Data.AcknowledgeTime.IsZero()

		if isAcknowledgedAlready {
			switch transfer.Data.AcknowledgeStatus {
			case models.TransferAccepted:
				mt += "\n" + whc.Translate(trans.MESSAGE_TEXT_ALREADY_ACCEPTED_TRANSFER)
			case models.TransferDeclined:
				mt += "\n" + whc.Translate(trans.MESSAGE_TEXT_ALREADY_DECLINED_TRANSFER)
			default:
				logus.Errorf(c, "!transfer.AcknowledgeTime.IsZero() && transfer.AcknowledgeStatus not in (accepted, declined)")
			}
		} else {
			mt += "\n" + whc.Translate(trans.MESSAGE_TEXT_PLEASE_ACKNOWLEDGE_TRANSFER)
		}
		receiptCode := receiptID

		if !isAcknowledgedAlready {
			inlineKeyboard = &tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{
						{
							Text:         whc.Translate(trans.COMMAND_TEXT_ACCEPT),
							CallbackData: fmt.Sprintf("%v?id=%v&do=%v", ACKNOWLEDGE_RECEIPT_CALLBACK_COMMAND, receiptCode, dtdal.AckAccept),
						},
					},
					{
						{
							Text:         whc.Translate(trans.COMMAND_TEXT_DECLINE),
							CallbackData: fmt.Sprintf("%v?id=%v&do=%v", ACKNOWLEDGE_RECEIPT_CALLBACK_COMMAND, receiptCode, dtdal.AckDecline),
						},
					},
				},
			}
		}
	}

	logus.Debugf(c, "mt: %v", mt)
	switch whc.InputType() {
	case botsfw.WebhookInputCallbackQuery:
		if m, err = whc.NewEditMessage(mt, botsfw.MessageFormatHTML); err != nil {
			return
		}
		m.DisableWebPagePreview = true
		if inlineKeyboard != nil {
			m.Keyboard = inlineKeyboard
		}
	case botsfw.WebhookInputText:
		m = whc.NewMessage(mt)
		if inlineKeyboard != nil {
			m.Keyboard = inlineKeyboard
		}
	default:
		if inputType, ok := botsfw.WebhookInputTypeNames[whc.InputType()]; ok {
			logus.Errorf(c, "Unknown input type: %d=%v", whc.InputType(), inputType)
		} else {
			logus.Errorf(c, "Unknown input type: %d", whc.InputType())
		}
	}

	if _, err = whc.Responder().SendMessage(c, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		if strings.Contains(err.Error(), "message is not modified") { // TODO: Can fail on different receipts for same amount
			logus.Warningf(c, fmt.Sprintf("Failed to send receipt to counterparty: %v", err))
		} else {
			return m, err
		}
	} else {
		if m, err = whc.NewEditMessage(
			whc.Translate(trans.MESSAGE_TEXT_RECEIPT_SENT_THROW_TELEGRAM)+"\n"+
				whc.Translate(trans.MESSAGE_TEXT_RECEIPT_VIEWED_BY_COUNTERPARTY),
			botsfw.MessageFormatHTML,
		); err != nil {
			return
		}
		m.EditMessageUID = telegram.NewChatMessageUID(transfer.Data.Creator().TgChatID, int(transfer.Data.CreatorTgReceiptByTgMsgID))
		//if _, err := whc.Responder().SendMessage(c, editCreatorMessage, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		//	logus.Errorf(c, "Failed to edit creator message: %v", err)
		//}
	}
	return m, err
}

func viewReceiptCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()

	logus.Debugf(c, "ViewReceiptAction(callbackUrl=%v)", callbackUrl)
	callbackQuery := callbackUrl.Query()

	localeCode5 := callbackQuery.Get("locale")
	if localeCode5 != "" {
		if err = whc.SetLocale(localeCode5); err != nil {
			return m, err
		}
		if _ /*appUser*/, err := whc.AppUserData(); err != nil {
			return m, err
		} else {
			panic("not implemented")
			//if _ = appUser.SetPreferredLocale(localeCode5); err != nil {
			//	return m, err
			//}
		}
	}
	receiptID := callbackQuery.Get("id")
	if receiptID == "" {
		return m, fmt.Errorf("receiptID is empty")
	}
	return ShowReceipt(whc, receiptID)
}

//func (viewReceiptCallback) onInvite(whc botsfw.WebhookContext, inviteCode string) (exit bool, transferID int, transfer *models.TransferEntry, m botsfw.MessageFromBot, err error) {
//	c := whc.Context()
//	var invite *invites.Invite
//	if invite, err = invites.GetInvite(c, inviteCode); err != nil {
//		return
//	} else {
//		if invite == nil {
//			err = fmt.Errorf("Invite not found by code: %v", inviteCode)
//			return
//		}
//		if invite.CreatedByUserID == whc.AppUserID() {
//			if transferID, err = invite.RelatedIntID(); err != nil {
//				return
//			}
//			if transfer, err = facade.Transfers.GetTransferByID(c, transferID); err != nil {
//				return
//			}
//			sender := whc.GetSender()
//			mt := getInlineReceiptMessage(whc, true, fmt.Sprintf("%v %v", sender.GetFirstName(), sender.GetLastName()))
//			editedMessage := tgbotapi.NewEditMessageTextByInlineMessageID(
//				whc.InputCallbackQuery().GetInlineMessageID(),
//				mt+"\n\n"+whc.Translate(trans.MESSAGE_TEXT_FOR_COUNTERPARTY_ONLY, transfer.ContactEntry().ContactName),
//			)
//			editedMessage.ParseMode = "HTML"
//			editedMessage.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
//				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
//					transferReceiptChooseLocaleButtons(inviteCode, invite.CreatedOnID, invite.CreatedOnPlatform),
//				},
//			}
//			m.TelegramEditMessageText = &editedMessage
//			exit = true
//			return
//		}
//
//		if transferID, transfer, _, _, err = ClaimInviteOnTransfer(whc, whc.InputCallbackQuery().GetInlineMessageID(), inviteCode, invite); err != nil {
//			return
//		}
//	}
//	return
//}
