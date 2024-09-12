package cmds4anybot

import (
	"github.com/bots-go-framework/bots-fw/botinput"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

const ChatLeftCommandCode = "left-chat"

var leftChatCommand = botsfw.Command{
	Code:       ChatLeftCommandCode,
	InputTypes: []botinput.WebhookInputType{botinput.WebhookInputLeftChatMembers},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return
	},
}
