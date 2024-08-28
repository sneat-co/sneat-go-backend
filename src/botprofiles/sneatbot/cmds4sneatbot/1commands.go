package cmds4sneatbot

import "github.com/bots-go-framework/bots-fw/botsfw"

// AddSneatSharedCommands registers commands shared by all Sneat bots
func AddSneatSharedCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	addCommands(commandsByType, []botsfw.Command{
		spacesCommand,
		membersCommand,
	})
}

// AddSneatBotCommands registers commands specific only to @SneatBot
func AddSneatBotCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command) {
	addCommands(commandsByType, []botsfw.Command{
		startCommand,
	})
}

// TODO: Decouple to shared package
func addCommands(commandsByType map[botsfw.WebhookInputType][]botsfw.Command, commands []botsfw.Command) {
	for _, c := range commands {
		for _, t := range c.InputTypes {
			commandsByType[t] = append(commandsByType[t], c)
		}
	}
}
