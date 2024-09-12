package shared_all

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func AddSharedRoutes(router botsfw.WebhooksRouter, botParams BotParams) {
	botParams.Validate()
	startCommand := createStartCommand(botParams.StartInBotAction, botParams.StartInGroupAction)
	helpRootCommand := createHelpRootCommand(botParams.HelpCommandAction)
	router.AddCommands(botinput.WebhookInputText, []botsfw.Command{
		startCommand,
		helpRootCommand,
		ReferrersCommand,
		createOnboardingAskLocaleCommand(botParams.SetMainMenu),
		aboutDrawCommand,
	})
	router.AddCommands(botinput.WebhookInputCallbackQuery, []botsfw.Command{
		onStartCallbackCommand(botParams.SetMainMenu),
		helpRootCommand,
		joinDrawCommand,
		aboutDrawCommand,
		AskPreferredLocaleFromSettingsCallback,
		setLocaleCallbackCommand(botParams.SetMainMenu),
	})
	router.AddCommands(botinput.WebhookInputLeftChatMembers, []botsfw.Command{
		leftChatCommand,
	})
	router.AddCommands(botinput.WebhookInputSticker, []botsfw.Command{
		botsfw.IgnoreCommand,
	})
	router.AddCommands(botinput.WebhookInputReferral, []botsfw.Command{
		startCommand,
	})
	router.AddCommands(botinput.WebhookInputConversationStarted, []botsfw.Command{
		startCommand,
	})
}
