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
	commandsByType := make(map[botinput.WebhookInputType][]botsfw.Command)
	cmds4anybot.AddSharedCommands(commandsByType)

	cmds4listusbot.AddListusOnlyBotCommands(commandsByType)
	cmds4listusbot.AddListusSharedCommands(commandsByType)
	router := botsfw.NewWebhookRouter(commandsByType, errFooterText)
	return anybot.NewProfile(ProfileID, &router)
}
