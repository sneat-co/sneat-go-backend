package sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/bots/listusbot"
	"github.com/sneat-co/sneat-go-backend/src/bots/shared"
)

const ProfileID = "sneat_bot"

var profile botsfw.BotProfile

func GetProfile(errFooterText func() string) botsfw.BotProfile {
	if profile == nil {
		profile = createProfile(errFooterText)
	}
	return profile
}

func createProfile(errFooterText func() string) botsfw.BotProfile {
	commandsByType := map[botsfw.WebhookInputType][]botsfw.Command{
		botsfw.WebhookInputText: []botsfw.Command{startCommand},
	}
	shared.AddSharedCommands(commandsByType)
	listusbot.AddListusCommands(commandsByType)
	router := botsfw.NewWebhookRouter(commandsByType, errFooterText)
	return shared.NewProfile(ProfileID, &router)
}
