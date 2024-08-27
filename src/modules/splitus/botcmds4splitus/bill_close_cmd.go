package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/logus"
	"net/url"
)

const CLOSE_BILL_COMMAND = "close-bill"

var closeBillCommand = billCallbackCommand(CLOSE_BILL_COMMAND,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		logus.Debugf(ctx, "closeBillCommand.CallbackAction()")
		return ShowBillCard(whc, true, bill, "Sorry, not implemented yet.")
	},
)
