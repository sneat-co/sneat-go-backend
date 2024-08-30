package cmds4listusbot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var startCommand = botsfw.Command{
	Code:     "start",
	Commands: []string{"/start"},
	Matcher: func(command botsfw.Command, context botsfw.WebhookContext) bool {
		return true
	},
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputText, botinput.WebhookInputInlineQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "Hello, stranger! I'm @Listus_bot. I can help you to manage your shopping list."
		return
	},
}
