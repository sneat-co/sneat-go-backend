package cmds4listusbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"github.com/sneat-co/sneat-go-backend/src/modules/listus/dbo4listus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"net/url"
)

func listCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	ctx := whc.Context()
	spaceRef := tghelpers.GetSpaceRef(callbackUrl)
	listKey := dbo4listus.ListKey(callbackUrl.Query().Get("k"))
	if err = listKey.Validate(); err != nil {
		return
	}
	if m, err = showListAction(ctx, whc, nil, spaceRef, listKey); err != nil {
		return
	}
	keyboard := m.Keyboard.(*tgbotapi.InlineKeyboardMarkup)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		tghelpers.BackToSpaceMenuButton(spaceRef),
	})

	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = keyboard

	m.ResponseChannel = botsfw.BotAPISendMessageOverHTTPS
	chatData := whc.ChatData()
	chatData.SetAwaitingReplyTo("list")
	switch chatData := chatData.(type) {
	case interface {
		SetSpaceRef(core4spaceus.SpaceRef)
	}:
		chatData.SetSpaceRef(spaceRef)
	default:
		err = fmt.Errorf("chatData %T does not support SetSpaceRef(core4spaceus.SpaceRef)", chatData)
	}
	return
}
