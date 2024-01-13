package common

import "github.com/bots-go-framework/bots-fw/botsfw"

var pingCommand = botsfw.Command{
	Code:       "ping",
	Commands:   []string{"/ping"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputText},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "Pong!"
		return
	},
}
