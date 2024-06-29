package splitus

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/strongo/logus"
	"net/url"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

const ADD_BILL_COMMENT_COMMAND = "bill_comment"

var addBillComment = billCallbackCommand(ADD_BILL_COMMENT_COMMAND,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		logus.Debugf(c, "addBillComment.CallbackAction()")

		//var editedMessage *tgbotapi.EditMessageTextConfig
		//if editedMessage, err = dtb_common.NewTelegramEditMessage(whc, "Enter new total for the bill:"); err != nil {
		//	return
		//}
		//
		//editedMessage.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, Selective: true}
		m = whc.NewMessage("Send your comment:")
		m.Keyboard = &tgbotapi.ForceReply{ForceReply: true, Selective: true}
		return
	},
)
