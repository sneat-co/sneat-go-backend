package tghelpers

import (
	telegram "github.com/bots-go-framework/bots-fw-telegram"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"strconv"
)

// GetEditMessageUID returns UID of the message to be edited
// TODO: Move to bots-fw-telegram?
func GetEditMessageUID(whc botsfw.WebhookContext) (*telegram.ChatMessageUID, error) {
	chatID, err := whc.Input().BotChatID()
	if err != nil {
		return nil, err
	}

	var tgChatID int64
	if tgChatID, err = strconv.ParseInt(chatID, 10, 64); err != nil {
		return nil, err
	}
	messageID := whc.Input().(telegram.TgWebhookCallbackQuery).GetMessage().IntID()
	return telegram.NewChatMessageUID(tgChatID, int(messageID)), nil
}
