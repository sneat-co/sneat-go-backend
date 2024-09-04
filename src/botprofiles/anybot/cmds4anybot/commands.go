package cmds4anybot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func AddSharedCommands(commandsByType map[botinput.WebhookInputType][]botsfw.Command) {
	for _, commands := range []botsfw.Command{
		pingCommand,
		counterCommand,
		contactMessageCommand,
	} {
		for _, inputType := range commands.InputTypes {
			commandsByType[inputType] = append(commandsByType[inputType], commands)
		}
	}
}
