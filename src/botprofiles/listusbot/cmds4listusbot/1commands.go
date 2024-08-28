package cmds4listusbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var Commands = []botsfw.Command{
	listCommand,
	addBuyItemCommand,
	remindCommand,
}

// AddListusBotCommands adds listus commands to a Listus bot
func AddListusBotCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	commandsByType[botsfw.WebhookInputText] = append(commandsByType[botsfw.WebhookInputText], startCommand)
	AddListusSharedCommands(commandsByType)
}

// AddListusSharedCommands adds anybot listus commands to a Sneat bot
func AddListusSharedCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	commandsByType[botsfw.WebhookInputText] = append(commandsByType[botsfw.WebhookInputText],
		listCommand,
		addBuyItemCommand,
		remindCommand,
	)
}
