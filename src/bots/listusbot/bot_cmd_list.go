package listusbot

import "github.com/bots-go-framework/bots-fw/botsfw"

var listCommand = botsfw.Command{
	Code:       "list",
	Commands:   []string{"/lisy"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputInlineQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "Hello, stranger!"
		return
	},
}
