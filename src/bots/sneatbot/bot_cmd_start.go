package sneatbot

import "github.com/bots-go-framework/bots-fw/botsfw"

var startCommand = botsfw.Command{
	Code:       "start",
	Commands:   []string{"/start"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputText, botsfw.WebhookInputInlineQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "Hello, stranger! I'm a @SneatBot. I can help you to manage your day-to-day family life."
		return
	},
}
