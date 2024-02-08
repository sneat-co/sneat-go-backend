package splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"net/url"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"github.com/strongo/log"
)

const CLOSE_BILL_COMMAND = "close-bill"

var closeBillCommand = billCallbackCommand(CLOSE_BILL_COMMAND,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		log.Debugf(c, "closeBillCommand.CallbackAction()")
		return ShowBillCard(whc, true, bill, "Sorry, not implemented yet.")
	},
)
