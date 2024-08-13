package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/facade4splitus"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/logus"
	"net/url"
)

var billChangeSplitModeCommand = botsfw.Command{
	Code: "split-mode",
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		logus.Debugf(c, "billChangeSplitModeCommand.CallbackAction()")
		var bill models4splitus.BillEntry
		if bill.ID, err = GetBillID(callbackUrl); err != nil {
			return
		}
		tx := whc.Tx()
		if bill, err = facade4splitus.GetBillByID(c, tx, bill.ID); err != nil {
			return
		}
		splitMode := models4splitus.SplitMode(callbackUrl.Query().Get("mode"))
		if bill.Data.SplitMode != splitMode {
			bill.Data.SplitMode = splitMode
			if err = facade4splitus.SaveBill(c, tx, bill); err != nil {
				return
			}
		}
		return ShowBillCard(whc, true, bill, "")
	},
}
