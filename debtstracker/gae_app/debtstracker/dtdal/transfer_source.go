package dtdal

import (
	"strconv"

	"github.com/bots-go-framework/bots-fw-telegram"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
)

type TransferSourceBot struct {
	platform string
	botID    string
	chatID   string
}

func (s TransferSourceBot) PopulateTransfer(t *models.TransferData) {
	t.CreatedOnPlatform = s.platform
	t.CreatedOnID = s.botID
	if s.platform == telegram.PlatformID {
		t.Creator().TgBotID = s.botID
		var err error
		t.Creator().TgChatID, err = strconv.ParseInt(s.chatID, 10, 64)
		if err != nil {
			panic(err.Error())
		}
	}
}

var _ TransferSource = (*TransferSourceBot)(nil)

func NewTransferSourceBot(platform, botID, chatID string) TransferSourceBot {
	if botID == "" {
		panic("botID is not set")
	}
	if chatID == "" {
		panic("chatID is not set")
	}
	return TransferSourceBot{
		platform: platform,
		botID:    botID,
		chatID:   chatID,
	}
}
