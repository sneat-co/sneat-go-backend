package dtb_transfer

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_general"
)

const CANCEL_TRANSFER_WIZARD_COMMAND = "cancel-transfer-wizard"

var CancelTransferWizardCommand = botsfw.Command{
	Code:     CANCEL_TRANSFER_WIZARD_COMMAND,
	Commands: trans.Commands(trans.COMMAND_TEXT_CANCEL, "/cancel", emoji.NO_ENTRY_SIGN_ICON),
	Action:   cancelTransferWizardCommandAction,
}

func cancelTransferWizardCommandAction(whc botsfw.WebhookContext) (botsfw.MessageFromBot, error) {
	whc.ChatData().SetAwaitingReplyTo("")
	//var m botsfw.MessageFromBot
	//userKey, _, err := whc.GetUser()
	//if err != nil {
	//	return m, err
	//}
	//var transfers []models.Transfer
	//ctx := whc.Context()
	//transferKeys, err := datastore.NewQuery(models.TransferKind).Filter("UserID =", userKey.IntID()).Limit(1).GetAll(ctx, &transfers)
	//if err != nil {
	//	return m, err
	//}
	m := whc.NewMessageByCode(trans.MESSAGE_TEXT_TRANSFER_CREATION_CANCELED)
	//if len(transferKeys) == 0 {
	//	m = tgbotapi.NewMessage(whc.ChatID(), Translate(trans.MESSAGE_TEXT_NOTHING_TO_CANCEL, whc))
	//} else {
	//	err := datastore.Delete(ctx, transferKeys[0])
	//	if err != nil {
	//		return m, err
	//	}
	//	//transfer := transfers[0]
	//}
	dtb_general.SetMainMenuKeyboard(whc, &m)
	return m, nil
}
