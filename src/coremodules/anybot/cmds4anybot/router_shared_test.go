package cmds4anybot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"testing"
)

func TestAddSharedAddSharedCommands(t *testing.T) {
	router := botsfw.NewWebhookRouter(nil)
	AddSharedCommands(router, BotParams{
		StartInBotAction: func(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error) {
			return
		},
		StartInGroupAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return
		},
		GetWelcomeMessageText: func(whc botsfw.WebhookContext) (text string, err error) {
			return
		},
		HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return
		},
		SetMainMenu: func(whc botsfw.WebhookContext, messageText string, showHint bool) (m botsfw.MessageFromBot, err error) {
			return
		},
	})
}
