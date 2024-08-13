package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/logus"
	"net/url"
)

const CHANGE_BILL_TOTAL_COMMAND = "bill_total"

var changeBillTotalCommand = billCallbackCommand(CHANGE_BILL_TOTAL_COMMAND,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		logus.Debugf(c, "changeBillTotalCommand.CallbackAction()")
		//var editedMessage *tgbotapi.EditMessageTextConfig
		//if editedMessage, err = dtb_common.NewTelegramEditMessage(whc, "Enter new total for the bill:"); err != nil {
		//	return
		//}
		//
		//editedMessage.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		m = whc.NewMessage("Enter new total for the bill:")
		m.Keyboard = &tgbotapi.ForceReply{ForceReply: true, Selective: true}
		return
	},
)
