package cmds4anybot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func AddSharedCommands(router botsfw.WebhooksRouter, botParams BotParams) {
	botParams.Validate()

	sharedCommands := []botsfw.Command{
		createStartCommand(
			botParams.StartInBotAction,
			botParams.StartInGroupAction,
			botParams.GetWelcomeMessageText,
			botParams.SetMainMenu,
		),
		spaceSettingsCommand,
		pingCommand,
		counterCommand,
		contactMessageCommand,
		createHelpRootCommand(botParams.HelpCommandAction, botParams.HelpCallbackAction),
		ReferrersCommand,
		UserSettingsLocaleCommand,
		leftChatCommand,
	}

	router.AddCommands(sharedCommands...)

	router.AddCommandsForInputType(botinput.WebhookInputSticker,
		botsfw.IgnoreCommand, // Can't add an input type to the command so must register explicitly
	)
}
