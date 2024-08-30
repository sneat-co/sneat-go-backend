package cmds4anybot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var pingCommand = botsfw.Command{
	Code:       "ping",
	Commands:   []string{"/ping"},
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputText},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "Pong!"
		return
	},
}
