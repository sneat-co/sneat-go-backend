package cmds4anybot

import "github.com/bots-go-framework/bots-fw/botsfw"

func AddSharedCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	for commandType, commands := range map[botsfw.WebhookInputType][]botsfw.Command{
		botsfw.WebhookInputText: {
			pingCommand,
			counterCommand,
		},
	} {
		commandsByType[commandType] = append(commandsByType[commandType], commands...)
	}
}
