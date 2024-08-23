package listusbot

import "github.com/bots-go-framework/bots-fw/botsfw"

var Commands = []botsfw.Command{
	listCommand,
	addBuyItemCommand,
	remindCommand,
}

// addListusBotCommands adds listus commands to a Listus bot
func addListusBotCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	commandsByType[botsfw.WebhookInputText] = append(commandsByType[botsfw.WebhookInputText], startCommand)
	AddListusCommands(commandsByType)
}

// AddListusCommands adds anybot listus commands to a Sneat bot
func AddListusCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	commandsByType[botsfw.WebhookInputText] = append(commandsByType[botsfw.WebhookInputText],
		listCommand,
		addBuyItemCommand,
		remindCommand,
	)
}
