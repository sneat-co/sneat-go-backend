package dtb_invite

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"net/url"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_transfer"
	"github.com/strongo/log"
)

var ChosenInlineResultCommand = botsfw.Command{
	Code:       "inline-create-invite",
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputChosenInlineResult},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		chosenResult := whc.Input().(botsfw.WebhookChosenInlineResult)
		query := chosenResult.GetQuery()
		log.Debugf(c, "ChosenInlineResultCommand.Action() => query: %v", query)

		queryUrl, err := url.Parse(query)
		if err != nil {
			return m, err
		}

		switch queryUrl.Path {
		case "receipt":
			return dtb_transfer.OnInlineChosenCreateReceipt(whc, chosenResult.GetInlineMessageID(), queryUrl)
		default:
			log.Warningf(c, "Unknown chosen inline query: "+query)
		}
		return
	},
}
