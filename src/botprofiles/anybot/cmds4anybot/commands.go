package cmds4anybot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func AddSharedCommands(commandsByType map[botinput.WebhookInputType][]botsfw.Command) {
	for commandType, commands := range map[botinput.WebhookInputType][]botsfw.Command{
		botinput.WebhookInputText: {
			pingCommand,
			counterCommand,
		},
	} {
		commandsByType[commandType] = append(commandsByType[commandType], commands...)
	}
}
