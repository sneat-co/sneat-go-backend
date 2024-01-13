package listusbot

import "github.com/bots-go-framework/bots-fw/botsfw"

var Commands = []botsfw.Command{
	listCommand,
	addBuyItemCommand,
	remindCommand,
}

var listusBotCommands = append(Commands,
	startCommand,
)
