package shared

import "github.com/bots-go-framework/bots-fw/botsfw"

func AddSharedCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	textCommands := commandsByType[botsfw.WebhookInputText]
	textCommands = append(textCommands, pingCommand)
	textCommands = append(textCommands, counterCommand)
	commandsByType[botsfw.WebhookInputText] = textCommands
}
