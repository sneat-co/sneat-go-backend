package splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/logus"
	"net/url"
)

const editBillCommandCode = "edit_bill"

var editBillCommand = billCallbackCommand(editBillCommandCode,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		logus.Debugf(c, "editBillCommand.CallbackAction()")
		var mt string

		if mt, err = getBillCardMessageText(c, whc.GetBotCode(), whc, bill, true, ""); err != nil {
			return
		}
		if m, err = whc.NewEditMessage(mt, botsfw.MessageFormatHTML); err != nil {
			return
		}
		m.Keyboard = getPrivateBillCardInlineKeyboard(whc, whc.GetBotCode(), bill)
		return
	},
)
