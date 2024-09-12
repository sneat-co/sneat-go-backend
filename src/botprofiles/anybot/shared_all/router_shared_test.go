package shared_all

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"testing"
)

func TestAddSharedRoutes(t *testing.T) {
	router := botsfw.NewWebhookRouter(map[botinput.WebhookInputType][]botsfw.Command{}, nil)
	AddSharedRoutes(router, BotParams{
		StartInBotAction: func(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error) {
			return
		},
		StartInGroupAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return
		},
		InBotWelcomeMessage: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return
		},
		HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return
		},
		SetMainMenu: func(whc botsfw.WebhookContext, m *botsfw.MessageFromBot) {

		},
	})
}
