package dtb_transfer

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/strongo/logus"
	"html"

	"errors"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_general"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/common"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
)

func AcknowledgeReceipt(whc botsfw.WebhookContext, receiptID, operation string) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()

	_, transfer, isCounterpartiesJustConnected, err := facade.AcknowledgeReceipt(c, receiptID, whc.AppUserID(), operation)
	if err != nil {
		if errors.Is(err, facade.ErrSelfAcknowledgement) {
			m = whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_SELF_ACKNOWLEDGEMENT, html.EscapeString(transfer.Data.Counterparty().ContactName)))
			return m, nil
		}
		return m, err
	} else {

		{ // Reporting to Google Analytics
			ga := whc.GA()

			if err = ga.Queue(ga.GaEventWithLabel(
				"receipts",
				"receipt-acknowledged",
				operation,
			)); err != nil {
				logus.Errorf(c, "Failed to report receipt-acknowledged to Google Analytics: %v", err)
			}

			if isCounterpartiesJustConnected {
				if err = ga.Queue(ga.GaEvent(
					"counterparties",
					"counterparties-connected",
				)); err != nil {
					logus.Errorf(c, "Failed to report counterparties-connected to Google Analytics: %v", err)
				}
			}
		}

		var operationMessage string
		switch operation {
		case dtdal.AckAccept:
			operationMessage = whc.Translate(trans.MESSAGE_TEXT_TRANSFER_ACCEPTED_BY_YOU)
		case dtdal.AckDecline:
			operationMessage = whc.Translate(trans.MESSAGE_TEXT_TRANSFER_DECLINED_BY_YOU)
		default:
			err = errors.New("Expected accept or decline as operation, got: " + operation)
			return
		}

		utm := common.NewUtmParams(whc, common.UTM_CAMPAIGN_RECEIPT)
		if whc.InputType() == botsfw.WebhookInputCallbackQuery {
			if m, err = whc.NewEditMessage(common.TextReceiptForTransfer(c, whc, transfer, "", common.ShowReceiptToCounterparty, utm)+"\n\n"+operationMessage, botsfw.MessageFormatHTML); err != nil {
				return
			}
		} else {
			m = whc.NewMessage(operationMessage + "\n\n" + common.TextReceiptForTransfer(c, whc, transfer, "", common.ShowReceiptToCounterparty, utm))
			m.Keyboard = dtb_general.MainMenuKeyboardOnReceiptAck(whc)
			m.Format = botsfw.MessageFormatHTML
		}

		if transfer.Data.Creator().TgChatID != 0 {
			askMsgToCreator := whc.NewMessage("")
			askMsgToCreator.ToChat = botsfw.ChatIntID(transfer.Data.Creator().TgChatID)
			var operationMsg string
			counterpartyName := transfer.Data.Counterparty().ContactName
			switch operation {
			case "accept":
				operationMsg = whc.Translate(trans.MESSAGE_TEXT_TRANSFER_ACCEPTED_BY_COUNTERPARTY, html.EscapeString(counterpartyName))
			case "decline":
				operationMsg = whc.Translate(trans.MESSAGE_TEXT_TRANSFER_DECLINED_BY_COUNTERPARTY, html.EscapeString(counterpartyName))
			default:
				err = errors.New("Expected accept or decline as operation, got: " + operation)
			}
			askMsgToCreator.Text = operationMsg + "\n\n" + common.TextReceiptForTransfer(c, whc, transfer, transfer.Data.CreatorUserID, common.ShowReceiptToAutodetect, utm)

			if transfer.Data.Creator().TgBotID != whc.GetBotCode() {
				logus.Warningf(c, "TODO: transferEntity.Creator().TgBotID != whc.GetBotCode(): "+askMsgToCreator.Text)
			} else {
				if _, err = whc.Responder().SendMessage(c, askMsgToCreator, botsfw.BotAPISendMessageOverHTTPS); err != nil {
					logus.Errorf(c, "Failed to send acknowledge to creator: %v", err)
					err = nil // This is not that critical to report the error to user
				}
			}
		}
		// Seems we can edit message just once after callback :(
		//if transferEntity.CounterpartyTgReceiptInlineMessageID != "" {
		//	mt = common.TextReceiptForTransfer(whc, transferID, transferEntity, transferEntity.CounterpartyCounterpartyID)
		//	editMessage := tgbotapi.NewEditMessageTextByInlineMessageID(transferEntity.CounterpartyTgReceiptInlineMessageID, mt + fmt.Sprintf("\n\n Acknowledged by %v", transferEntity.ContactEntry().ContactName))
		//
		//	if values, err := editMessage.Values(); err != nil {
		//		logus.Errorf(c, "Failed to get values for editMessage: %v", err)
		//	} else {
		//		logus.Debugf(c, "editMessage.Values(): %v", values)
		//	}
		//	updateMessage := whc.NewMessage("")
		//	updateMessage.TelegramEditMessageText = &editMessage
		//	_, err := whc.Responder().SendMessage(c, updateMessage, botsfw.BotAPISendMessageOverHTTPS)
		//	if err != nil {
		//		logus.Errorf(c, "Failed to update counterparty receipt message: %v", err)
		//	}
		//}
		return m, err
	}
}
