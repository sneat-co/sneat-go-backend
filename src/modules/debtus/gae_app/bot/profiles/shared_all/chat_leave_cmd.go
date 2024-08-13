package shared_all

import "github.com/bots-go-framework/bots-fw/botsfw"

const CHAT_LEFT_COMMAND = "left-chat"

var leftChatCommand = botsfw.Command{
	Code: CHAT_LEFT_COMMAND,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		return
	},
}
