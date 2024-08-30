package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
)

var startCommand = botsfw.Command{
	Code:     "start",
	Commands: []string{"/start"},
	InputTypes: []botinput.WebhookInputType{
		botinput.WebhookInputText,
		botinput.WebhookInputCallbackQuery,
		botinput.WebhookInputInlineQuery,
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
	m, err = spaceAction(whc, core4spaceus.NewSpaceRef(core4spaceus.SpaceTypeFamily, ""))
	return
}
