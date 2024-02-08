package dtb_transfer

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"strings"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/dtdal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/facade"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/general"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/invites"
	"github.com/strongo/log"
)

const ASK_EMAIL_FOR_RECEIPT_COMMAND = "ask-email-for-receipt"

var AskEmailForReceiptCommand = botsfw.Command{
	Code: ASK_EMAIL_FOR_RECEIPT_COMMAND,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()

		log.Debugf(c, "AskEmailForReceiptCommand.Action()")
		email := whc.Input().(botsfw.WebhookTextMessage).Text()
		if !strings.Contains(email, "@") {
			return whc.NewMessage(whc.Translate(trans.MESSAGE_TEXT_INVALID_EMAIL)), nil
		}

		chatEntity := whc.ChatData()
		transferID := chatEntity.GetWizardParam(WIZARD_PARAM_TRANSFER)
		transfer, err := facade.Transfers.GetTransferByID(c, nil, transferID)
		if err != nil {
			return m, err
		}
		m, err = sendReceiptByEmail(whc, email, "", transfer)
		return
	},
}

func sendReceiptByEmail(whc botsfw.WebhookContext, toEmail, toName string, transfer models.Transfer) (m botsfw.MessageFromBot, err error) {
	c := whc.Context()
	receiptEntity := models.NewReceiptEntity(whc.AppUserID(), transfer.ID, transfer.Data.Counterparty().UserID, whc.Locale().Code5, string(models.InviteByEmail), toEmail, general.CreatedOn{
		CreatedOnPlatform: whc.BotPlatform().ID(),
		CreatedOnID:       whc.GetBotCode(),
	})
	var receipt models.Receipt
	if receipt, err = dtdal.Receipt.CreateReceipt(c, receiptEntity); err != nil {
		return m, err
	}

	emailID := ""
	if emailID, err = invites.SendReceiptByEmail(
		c,
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
