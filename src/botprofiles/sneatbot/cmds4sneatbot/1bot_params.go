package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/cmds4anybot"
)

func GetBotParams() cmds4anybot.BotParams {
	return cmds4anybot.BotParams{
		GetWelcomeMessageText: sneatBotWelcomeMessage,
		StartInBotAction:      startActionWithStartParams,
		StartInGroupAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			m.Text = "Start in group is not implemented yet for @SneatBot"
			return
		},
		HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			m.Text = "Help is not implemented yet for @SneatBot"
			return
		},
		SetMainMenu: func(whc botsfw.WebhookContext, m *botsfw.MessageFromBot) {
			m.Keyboard = cmds4anybot.StartMessageInlineKeyboard(whc)
		},
	}
}
