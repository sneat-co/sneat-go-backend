package debtus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
	"strings"

	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_inline"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_invite"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_transfer"
)

var InlineQueryCommand = botsfw.Command{
	InputTypes: nil,
	Icon:       "",
	Replies:    nil,
	Code:       "inline-query",
	Title:      "",
	Titles:     nil,
	ExactMatch: "",
	Commands:   nil,
	Matcher:    nil,
	//InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputInlineQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()

		inlineQuery := whc.Input().(botsfw.WebhookInlineQuery)
		query := inlineQuery.GetQuery()
		logus.Debugf(c, "InlineQueryCommand.Action(query=%v)", query)
		switch {
		case query == "":
			m, err = dtb_inline.InlineEmptyQuery(whc)
		case query == "/invite":
			m, err = dtb_invite.InlineSendInvite(whc)
		case strings.HasPrefix(query, "receipt?id="):
			m, err = dtb_transfer.InlineSendReceipt(whc)
		//case strings.HasPrefix(query, "accept?transfer="):
		//	m, err = dtb_transfer.InlineAcceptTransfer(whc)
		default:
			amountMatches := dtb_inline.ReInlineQueryAmount.FindStringSubmatch(query)
			if amountMatches != nil {
				return dtb_inline.InlineNewRecord(whc, amountMatches)
			}
			logus.Debugf(c, "Inline query not matched to any action: [%v]", query)
		}
		return
	},
	CallbackAction: nil,
}
