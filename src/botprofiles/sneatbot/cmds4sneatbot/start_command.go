package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var startCommand = botsfw.Command{
	Code:     "start",
	Commands: []string{"/start"},
	InputTypes: []botsfw.WebhookInputType{
		botsfw.WebhookInputText,
		botsfw.WebhookInputCallbackQuery,
		botsfw.WebhookInputInlineQuery,
	},
	Action: startAction,
}

func startAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {

	var welcomeMsg botsfw.MessageFromBot
	welcomeMsg.Text = "Hello, stranger!"
	welcomeMsg.Text += "\n\nI'm a @SneatBot. I can help you to manage your day-to-day family life."
	welcomeMsg.Text += "\n\nOr you can create a space to manage your group/team/community."

	responder := whc.Responder()
	ctx := whc.Context()
	if _, err = responder.SendMessage(ctx, welcomeMsg, botsfw.BotAPISendMessageOverHTTPS); err != nil {
		return
	}
	m, err = spaceAction(whc, "family")
	return
}
