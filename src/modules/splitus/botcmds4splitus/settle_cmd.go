package botcmds4splitus

import (
	"errors"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"net/url"
)

const settleBillsCommandCode = "settle"

var settleBillsCommand = botsfw.Command{
	Code:     settleBillsCommandCode,
	Commands: []string{"/settle"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return settleBillsAction(whc)
	},
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		return settleBillsAction(whc)
	},
}

func settleBillsAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	logus.Debugf(ctx, "settleBillsAction()")
	err = errors.New("settleBillsAction not implemented yet")
	//var user models4debtus.AppUser
	//if user, err = dal4userus.GetUserByID(ctx, nil, whc.AppUserID()); err != nil {
	//	return
	//}
	//
	//outstandingBills := user.Data.BillsHolder.GetOutstandingBills()
	//
	//m.Text = fmt.Sprintf("len(outstandingBills): %v", len(outstandingBills))

	return
}
