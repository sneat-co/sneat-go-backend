package cmds4sneatbot

import "github.com/bots-go-framework/bots-fw/botsfw"

// AddSneatSharedCommands registers commands shared by all Sneat bots
func AddSneatSharedCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	for _, c := range []botsfw.Command{
		spacesCommand,
		membersCommand,
	} {
		for _, t := range c.InputTypes {
			commandsByType[t] = append(commandsByType[t], c)
		}
	}
}

// AddSneatBotCommands registers commands specific only to @SneatBot
func AddSneatBotCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	commandsByType[botsfw.WebhookInputText] = append(commandsByType[botsfw.WebhookInputText],
		startCommand,
	)
}
