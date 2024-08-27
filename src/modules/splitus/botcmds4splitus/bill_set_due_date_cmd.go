package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"net/url"
)

const setBillDueDateCommandCode = "bill_due"

var setBillDueDateCommand = botsfw.Command{
	Code: setBillDueDateCommandCode,
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		chatEntity := whc.ChatData()
		chatEntity.SetAwaitingReplyTo(setBillDueDateCommandCode)
		chatEntity.AddWizardParam("bill", callbackUrl.Query().Get("id"))
		logus.Debugf(ctx, "setBillDueDateCommand.CallbackAction()")
		m = whc.NewMessage("Please set bill due date as dd.mm.yyyy")
		m.Keyboard = &tgbotapi.ForceReply{ForceReply: true, Selective: true}
		return
	},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		logus.Debugf(ctx, "setBillDueDateCommand.Action()")
		m = whc.NewMessage("Not implemented yet")
		return
	},
}
