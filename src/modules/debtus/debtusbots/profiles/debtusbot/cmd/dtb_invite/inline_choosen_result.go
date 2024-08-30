package dtb_invite

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/debtusbot/cmd/dtb_transfer"
	"github.com/strongo/logus"
	"net/url"
)

var ChosenInlineResultCommand = botsfw.Command{
	Code:       "inline-create-invite",
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputChosenInlineResult},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		chosenResult := whc.Input().(botinput.WebhookChosenInlineResult)
		query := chosenResult.GetQuery()
		logus.Debugf(ctx, "ChosenInlineResultCommand.Action() => query: %v", query)

		queryUrl, err := url.Parse(query)
		if err != nil {
			return m, err
		}

		switch queryUrl.Path {
		case "receipt":
			return dtb_transfer.OnInlineChosenCreateReceipt(whc, chosenResult.GetInlineMessageID(), queryUrl)
		default:
			logus.Warningf(ctx, "Unknown chosen inline query: "+query)
		}
		return
	},
}
