package debtus

import (
	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/log"
)

const NEW_CHAT_MEMBERS_COMMAND = "new-chat-members"

var newChatMembersCommand = botsfw.Command{
	Code: NEW_CHAT_MEMBERS_COMMAND,
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		if whc.IsInGroup() {
			log.Warningf(whc.Context(), "Leaving chat as @DebtsTrackerBot does not support group chats")
			m.BotMessage = telegram.LeaveChat{}
		}
		return
	},
}
