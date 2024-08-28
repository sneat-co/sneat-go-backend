package sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/listusbot/cmds4listusbot"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/sneatbot/cmds4sneatbot"
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
	commandsByType := make(map[botsfw.WebhookInputType][]botsfw.Command)
	cmds4anybot.AddSharedCommands(commandsByType)
	cmds4sneatbot.AddSneatBotCommands(commandsByType)
	cmds4sneatbot.AddSneatSharedCommands(commandsByType)
	cmds4listusbot.AddListusSharedCommands(commandsByType)
	router := botsfw.NewWebhookRouter(commandsByType, errFooterText)
	return anybot.NewProfile(ProfileID, &router)
}
