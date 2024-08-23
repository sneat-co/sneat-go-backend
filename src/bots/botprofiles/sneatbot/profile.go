package sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	shared2 "github.com/sneat-co/sneat-go-backend/src/bots/botprofiles/anybot"
	"github.com/sneat-co/sneat-go-backend/src/bots/botprofiles/listusbot"
)

const ProfileID = "sneat_bot"

var profile botsfw.BotProfile

func GetProfile(errFooterText func() string) botsfw.BotProfile {
	if profile == nil {
		profile = createSneatBotProfile(errFooterText)
	}
	return profile
}

func createSneatBotProfile(errFooterText func() string) botsfw.BotProfile {
	commandsByType := map[botsfw.WebhookInputType][]botsfw.Command{
		botsfw.WebhookInputText: {startCommand},
	}
	shared2.AddSharedCommands(commandsByType)
	listusbot.AddListusCommands(commandsByType)
	router := botsfw.NewWebhookRouter(commandsByType, errFooterText)
	return shared2.NewProfile(ProfileID, &router)
}
