package sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/bots/common"
	"github.com/sneat-co/sneat-go-backend/src/bots/listusbot"
)

const ProfileID = "sneat_bot"

var Profile botsfw.BotProfile

func init() {
	var textAndContactCommands = []botsfw.Command{startCommand}
	textAndContactCommands = append(textAndContactCommands, common.Commands...)
	textAndContactCommands = append(textAndContactCommands, listusbot.Commands...)

	commandsByType := map[botsfw.WebhookInputType][]botsfw.Command{
		botsfw.WebhookInputText: textAndContactCommands,
	}
	router := botsfw.NewWebhookRouter(commandsByType, func() string {
		return "Please report any issues to @trakhimenok"
	})
	Profile = common.NewProfile(ProfileID, &router)
}
