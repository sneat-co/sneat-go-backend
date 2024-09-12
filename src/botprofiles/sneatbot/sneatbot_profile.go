package sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/anybot/cmds4anybot"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/listusbot/cmds4listusbot"
	"github.com/sneat-co/sneat-go-backend/src/botprofiles/sneatbot/cmds4sneatbot"
)

const ProfileID = "sneat_bot"

var sneatBotProfile botsfw.BotProfile

func GetProfile(errFooterText func() string) botsfw.BotProfile {
	if sneatBotProfile == nil {
		sneatBotProfile = createSneatBotProfile(errFooterText)
	}
	return sneatBotProfile
}

func createSneatBotProfile(errFooterText func() string) botsfw.BotProfile {
	router := botsfw.NewWebhookRouter(errFooterText)

	botParams := cmds4sneatbot.GetBotParams()
	cmds4anybot.AddSharedCommands(router, botParams)

	commandsByType := make(map[botinput.WebhookInputType][]botsfw.Command) // TODO: get rid of `commandsByType`
	//cmds4sneatbot.AddSneatBotCommands(commandsByType)
	cmds4sneatbot.AddSneatSharedCommands(commandsByType)
	cmds4listusbot.AddListusSharedCommands(commandsByType)
	return anybot.NewProfile(ProfileID, &router)
}
