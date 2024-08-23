package listusbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	shared2 "github.com/sneat-co/sneat-go-backend/src/bots/botprofiles/anybot"
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
	commandsByType := map[botsfw.WebhookInputType][]botsfw.Command{
		botsfw.WebhookInputText: []botsfw.Command{startCommand},
	}
	shared2.AddSharedCommands(commandsByType)
	addListusBotCommands(commandsByType)
	router := botsfw.NewWebhookRouter(commandsByType, errFooterText)
	return shared2.NewProfile(ProfileID, &router)
}
