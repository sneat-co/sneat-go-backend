package collectus

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-core-modules/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-core-modules/userus/dbo4userus"
)

var botParams = cmds4anybot.BotParams{
	StartInBotAction: func(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error) {
		m.Text = "StartInBotAction is not implemented yet"
		return
	},
	StartInGroupAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "StartInGroupAction is not implemented yet"
		return
	},
	HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "HelpCommandAction is not implemented yet"
		return
	},
	SetMainMenu: func(whc botsfw.WebhookContext, messageText string, showHint bool) (m botsfw.MessageFromBot, err error) {
		m.Text = "Collectus main menu"
		return
	},
	GetWelcomeMessageText: func(whc botsfw.WebhookContext) (text string, err error) {
		var user dbo4userus.UserEntry
		if user, err = cmds4anybot.GetUser(whc); err != nil {
			return
		}
		text = whc.Translate(
			trans.MESSAGE_TEXT_HI_USERNAME, user.Data.Names.FirstName) + " " + whc.Translate(trans.COLLECTUS_TEXT_HI) +
			"\n\n" + whc.Translate(trans.COLLECTUS_TEXT_ABOUT_ME_AND_CO) +
			"\n\n" + whc.Translate(trans.COLLECTUS_TG_COMMANDS)

		//m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		//	[]tgbotapi.InlineKeyboardButton{
		//		tgbotapi.NewInlineKeyboardButtonSwitchInlineQuery(
		//			whc.CommandText(trans.COMMAND_TEXT_NEW_FUNDRAISING, emoji.MEMO_ICON),
		//			"",
		//		),
		//	},
		//	//[]tgbotapi.InlineKeyboardButton{
		//	//	shared_all.NewGroupTelegramInlineButton(whc, 0),
		//	//},
		//)
		return
	},
}

var Router = botsfw.NewWebhookRouter(
	func() string { return "Please report any errors to @CollectusGroup" },
)

func init() {
	cmds4anybot.AddSharedCommands(Router, botParams)
}
