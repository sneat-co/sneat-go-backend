package listusbot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/listusbot/cmds4listusbot"
)

const ProfileID = "listus_bot"

var profile botsfw.BotProfile

func GetProfile(errFooterText func() string) botsfw.BotProfile {
	if profile == nil {
		profile = createProfile(errFooterText)
	}
	return profile
}

func createProfile(errFooterText func() string) botsfw.BotProfile {
	botParams := cmds4anybot.BotParams{
		StartInBotAction: func(whc botsfw.WebhookContext, startParams []string) (m botsfw.MessageFromBot, err error) {
			return
		},
		StartInGroupAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return
		},
		GetWelcomeMessageText: func(whc botsfw.WebhookContext) (text string, err error) {
			return "Welcome to Listus bot", nil
		},
		HelpCommandAction: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
			return
		},
		SetMainMenu: func(whc botsfw.WebhookContext, m *botsfw.MessageFromBot) {
		},
	}
	router := botsfw.NewWebhookRouter(errFooterText)
	cmds4anybot.AddSharedCommands(router, botParams)

	commandsByType := make(map[botinput.WebhookInputType][]botsfw.Command)
	cmds4listusbot.AddListusSharedCommands(commandsByType)
	router.AddCommandsGroupedByType(commandsByType)
	return anybot.NewProfile(ProfileID, &router)
}
