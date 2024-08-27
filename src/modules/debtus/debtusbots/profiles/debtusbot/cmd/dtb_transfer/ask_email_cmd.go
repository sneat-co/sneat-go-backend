package dtb_transfer

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/facade4debtus"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/general"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/gae_app/invites"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/models4debtus"
	"github.com/strongo/logus"
	"strings"
)

const ASK_EMAIL_FOR_RECEIPT_COMMAND = "ask-email-for-receipt"

var AskEmailForReceiptCommand = botsfw.Command{
	Code: ASK_EMAIL_FOR_RECEIPT_COMMAND,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()

		logus.Debugf(ctx, "AskEmailForReceiptCommand.Action()")
		email := whc.Input().(botsfw.WebhookTextMessage).Text()
		if !strings.Contains(email, "@") {
			return whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_INVALID_EMAIL)), nil
		}

		chatEntity := whc.ChatData()
		transferID := chatEntity.GetWizardParam(WizardParamTransfer)
		transfer, err := facade4debtus.Transfers.GetTransferByID(ctx, nil, transferID)
		if err != nil {
			return m, err
		}
		m, err = sendReceiptByEmail(whc, email, "", transfer)
		return
	},
}

func sendReceiptByEmail(whc botsfw.WebhookContext, toEmail, toName string, transfer models4debtus.TransferEntry) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	receiptEntity := models4debtus.NewReceiptEntity(whc.AppUserID(), transfer.ID, transfer.Data.Counterparty().UserID, whc.Locale().Code5, string(models4debtus.InviteByEmail), toEmail, general.CreatedOn{
		CreatedOnPlatform: whc.BotPlatform().ID(),
		CreatedOnID:       whc.GetBotCode(),
	})
	var receipt models4debtus.ReceiptEntry
	if receipt, err = dtdal.Receipt.CreateReceipt(ctx, receiptEntity); err != nil {
		return m, err
	}

	emailID := ""
	if emailID, err = invites.SendReceiptByEmail(
		ctx,
		whc,
		receipt,
		whc.GetSender().GetFirstName(),
		toName,
		toEmail,
	); err != nil {
		return m, err
	}

	m = whc.NewMessageByCode(trans.MESSAGE_TEXT_RECEIPT_SENT_THROW_EMAIL, emailID)

	return m, err
}
