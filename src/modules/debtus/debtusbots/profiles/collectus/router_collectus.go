package collectus

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
	shared_all2 "github.com/sneat-co/sneat-go-backend/src/modules/debtus/debtusbots/profiles/shared_all"
	"github.com/sneat-co/sneat-go-backend/src/modules/userus/dbo4userus"
)

var botParams = shared_all2.BotParams{
	InBotWelcomeMessage: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		var user dbo4userus.UserEntry
		if user, err = shared_all2.GetUser(whc); err != nil {
			return
		}
		m.Text = whc.Translate(
			trans.MESSAGE_TEXT_HI_USERNAME, user.Data.Names.FirstName) + " " + whc.Translate(trans.COLLECTUS_TEXT_HI) +
			"\n\n" + whc.Translate(trans.COLLECTUS_TEXT_ABOUT_ME_AND_CO) +
			"\n\n" + whc.Translate(trans.COLLECTUS_TG_COMMANDS)
		m.Format = botsfw.MessageFormatHTML
		m.IsEdit = true

		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
					whc.CommandText(trans.COMMAND_TEXT_NEW_FUNDRAISING, emoji.MEMO_ICON),
					"",
				),
			},
			//[]tgbotapi.InlineKeyboardButton{
			//	shared_all.NewGroupTelegramInlineButton(whc, 0),
			//},
		)
		return
	},
}

var Router = botsfw.NewWebhookRouter(
	map[botinput.WebhookInputType][]botsfw.Command{},
	func() string { return "Please report any errors to @CollectusGroup" },
)

func init() {
	shared_all2.AddSharedRoutes(Router, botParams)
}
