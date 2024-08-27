package dtb_transfer

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/common4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dal4userus"
	"github.com/strongo/logus"
	"html"
	"net/url"

	"github.com/bots-go-framework/bots-fw-telegram"
	"strings"
)

//func CancelReceiptAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
//	return whc.NewMessage("TODO: Sorry, cancel is not implemented yet..."), nil
//}

const VIEW_RECEIPT_CALLBACK_COMMAND = "view-receipt"

var ViewReceiptCallbackCommand = botsfw.NewCallbackCommand(VIEW_RECEIPT_CALLBACK_COMMAND, viewReceiptCallbackAction)

func ShowReceipt(whc botsfw.WebhookContext, receiptID string) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	var receipt models4debtus.ReceiptEntry
	if receipt, err = dtdal.Receipt.GetReceiptByID(ctx, nil, receiptID); err != nil {
		return m, err
	}

	if receipt.Data.CreatorUserID == whc.AppUserID() {
		m.Text = whc.Translate(trans.MESSAGE_TEXT_RECEIPT_ATTEMPT_TO_VIEW_OWN)
		return
	}

	receipt, err = facade4debtus.MarkReceiptAsViewed(ctx, receiptID, whc.AppUserID())
	if err != nil {
		return
	}

	transfer, err := facade4debtus.Transfers.GetTransferByID(ctx, nil, receipt.Data.TransferID)
	if err != nil {
		return m, err
	}

	m = whc.NewMessage("")

	var (
		mt           string
		counterparty models4debtus.DebtusSpaceContactEntry
	)
	counterpartyCounterparty := transfer.Data.Creator()

	if counterpartyCounterparty.ContactID != "" {
		counterparty, err = facade4debtus.GetDebtusSpaceContactByID(ctx, nil, receipt.Data.SpaceID, counterpartyCounterparty.ContactID)
	} else {
		if user, err := dal4userus.GetUserByID(ctx, nil, transfer.Data.CreatorUserID); err != nil {
			return m, err
		} else {
			counterparty.Data = &models4debtus.DebtusSpaceContactDbo{}
			counterparty.Data.FirstName = user.Data.Names.FirstName
			counterparty.Data.LastName = user.Data.Names.LastName
		}
	}

	if err != nil {
		return m, err
	}
	utm := common4debtus.NewUtmParams(whc, common4debtus.UTM_CAMPAIGN_REMINDER)
	mt = common4debtus.TextReceiptForTransfer(ctx, whc, transfer, whc.AppUserID(), common4debtus.ShowReceiptToAutodetect, utm)

	logus.Debugf(ctx, "ReceiptEntry text: %v", mt)

	var inlineKeyboard *tgbotapi.InlineKeyboardMarkup

	if receipt.Data.CreatorUserID == whc.AppUserID() {
		mt += "\n" + whc.Translate(trans.MESSAGE_TEXT_SELF_ACKNOWLEDGEMENT, html.EscapeString(transfer.Data.Counterparty().ContactName))
	} else {
		isAcknowledgedAlready := !transfer.Data.AcknowledgeTime.IsZero()

		if isAcknowledgedAlready {
			switch transfer.Data.AcknowledgeStatus {
			case models4debtus.TransferAccepted:
				mt += "\n" + whc.Translate(trans.MESSAGE_TEXT_ALREADY_ACCEPTED_TRANSFER)
			case models4debtus.TransferDeclined:
				mt += "\n" + whc.Translate(trans.MESSAGE_TEXT_ALREADY_DECLINED_TRANSFER)
			default:
				logus.Errorf(ctx, "!transfer.AcknowledgeTime.IsZero() && transfer.AcknowledgeStatus not in (accepted, declined)")
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

	logus.Debugf(ctx, "mt: %v", mt)
	switch inputType := whc.InputType(); inputType {
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
		logus.Errorf(ctx, "Unknown input type: %s", botsfw.GetWebhookInputTypeIdNameString(inputType))
	}

	if _, err = whc.Responder().SendMessage(ctx, m, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		if strings.Contains(err.Error(), "message is not modified") { // TODO: Can fail on different receipts for same amount
			logus.Warningf(ctx, fmt.Sprintf("Failed to send receipt to counterparty: %v", err))
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
		//if _, err := whc.Responder().SendMessage(ctx, editCreatorMessage, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		//	logus.Errorf(ctx, "Failed to edit creator message: %v", err)
		//}
	}
	return m, err
}

func viewReceiptCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()

	logus.Debugf(ctx, "ViewReceiptAction(callbackUrl=%v)", callbackUrl)
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
//	ctx := whc.Context()
//	var invite *invites.Invite
//	if invite, err = invites.GetInvite(ctx, inviteCode); err != nil {
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
//			if transfer, err = facade4debtus.QueueTransfers.GetTransferByID(c, transferID); err != nil {
//				return
//			}
//			sender := whc.GetSender()
//			mt := getInlineReceiptMessage(whc, true, fmt.Sprintf("%v %v", sender.GetFirstName(), sender.GetLastName()))
//			editedMessage := tgbotapi.NewEditMessageTextByInlineMessageID(
//				whc.InputCallbackQuery().GetInlineMessageID(),
//				mt+"\n\n"+whc.Translate(trans.MESSAGE_TEXT_FOR_COUNTERPARTY_ONLY, transfer.DebtusSpaceContactEntry().ContactName),
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
