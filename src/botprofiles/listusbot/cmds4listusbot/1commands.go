package cmds4listusbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/bothelpers"
)

// AddListusOnlyBotCommands adds listus commands to a Listus bot
func AddListusOnlyBotCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	bothelpers.AddCommands(commandsByType, []botsfw.Command{
		startCommand,
	})
}

// AddListusSharedCommands adds listus commands to a Sneat bot
func AddListusSharedCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	bothelpers.AddCommands(commandsByType, []botsfw.Command{
		listCommand,
		todoCommand,
		watchCommand,
		remindCommand,
	})
}
