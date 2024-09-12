package cmds4anybot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func AddSharedCommands(router botsfw.WebhooksRouter, botParams BotParams) {
	botParams.Validate()

	sharedCommands := []botsfw.Command{
		createStartCommand(botParams.StartInBotAction, botParams.StartInGroupAction, botParams.GetWelcomeMessageText),
		newStartCallbackCommand(botParams.SetMainMenu), // TODO: Should be part of the start command
		onboardingCommand,
		pingCommand,
		counterCommand,
		contactMessageCommand,
		createHelpRootCommand(botParams.HelpCommandAction),
		ReferrersCommand,
		createOnboardingAskLocaleCommand(botParams.SetMainMenu),
		aboutDrawCommand, joinDrawCommand, // should be in a dedicated module and be registered at once
		AskPreferredLocaleFromSettingsCallback,
		newSetLocaleCallbackCommand(botParams.SetMainMenu),
		leftChatCommand,
	}

	router.AddCommands(sharedCommands...)

	router.AddCommandsForInputType(botinput.WebhookInputSticker,
		botsfw.IgnoreCommand, // Can't add an input type to the command so must register explicitly
	)
}
