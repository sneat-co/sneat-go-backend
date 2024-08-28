package bothelpers

import "github.com/bots-go-framework/bots-fw/botsfw"

func AddCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command, commands []botsfw.Command) {
	for _, c := range commands {
		for _, t := range c.InputTypes {
			commandsByType[t] = append(commandsByType[t], c)
		}
	}
}
