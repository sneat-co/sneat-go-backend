package cmds4listusbot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-core-modules/botscore/bothelpers"
)

// AddListusSharedCommands adds listus commands to a Sneat bot
func AddListusSharedCommands(commandsByType map[botinput.WebhookInputType][]botsfw.Command) {
	bothelpers.AddCommands(commandsByType, []botsfw.Command{
		listusListCommand,
		todoCommand,
		watchCommand,
		remindCommand,
	})
}
