package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/bothelpers"
)

// AddSneatSharedCommands registers commands shared by all Sneat bots
func AddSneatSharedCommands(commandsByType map[botinput.WebhookInputType][]botsfw.Command) {
	bothelpers.AddCommands(commandsByType, []botsfw.Command{
		spacesCommand,
		spaceCommand,
		membersCommand,
		contactsCommand,
		assetsCommand,
		budgetCommand,
		debtsCommand,
		settingsCommand,
		calendarCommand,
	})
}

// AddSneatBotCommands registers commands specific only to @SneatBot
//func AddSneatBotCommands(commandsByType map[botinput.WebhookInputType][]botsfw.Command) {
//	bothelpers.AddCommands(commandsByType, []botsfw.Command{
//		startCommand,
//	})
//}
