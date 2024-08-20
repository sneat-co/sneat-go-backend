package shared_all

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
)

func AddSharedRoutes(router botsfw.WebhooksRouter, botParams BotParams) {
	startCommand := createStartCommand(botParams)
	helpRootCommand := createHelpRootCommand(botParams)
	router.AddCommands(botsfw.WebhookInputText, []botsfw.Command{
		startCommand,
		helpRootCommand,
		ReferrersCommand,
		createOnboardingAskLocaleCommand(botParams),
		aboutDrawCommand,
	})
	router.AddCommands(botsfw.WebhookInputCallbackQuery, []botsfw.Command{
		onStartCallbackCommand(botParams),
		helpRootCommand,
		joinDrawCommand,
		aboutDrawCommand,
		askPreferredLocaleFromSettingsCallback,
		setLocaleCallbackCommand(botParams),
	})
	router.AddCommands(botsfw.WebhookInputLeftChatMembers, []botsfw.Command{
		leftChatCommand,
	})
	router.AddCommands(botsfw.WebhookInputSticker, []botsfw.Command{
		botsfw.IgnoreCommand,
	})
	router.AddCommands(botsfw.WebhookInputReferral, []botsfw.Command{
		startCommand,
	})
	router.AddCommands(botsfw.WebhookInputConversationStarted, []botsfw.Command{
		startCommand,
	})
}
