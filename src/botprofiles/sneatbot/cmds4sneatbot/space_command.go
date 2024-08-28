package cmds4sneatbot

import (
	"fmt"
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"net/url"
	"strings"
)

var spaceCommand = botsfw.Command{
	Code:           "space",
	Commands:       []string{"/space"},
	InputTypes:     []botsfw.WebhookInputType{botsfw.WebhookInputCallbackQuery},
	CallbackAction: spaceCallbackAction,
}

func spaceCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	spaceID := callbackUrl.Query().Get("id")
	if m, err = spaceAction(whc, spaceID); err != nil {
		return
	}
	keyboard := m.Keyboard
	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = keyboard
	if m.EditMessageUID, err = tghelpers.GetEditMessageUID(whc); err != nil {
		return
	}
	return
}

func spaceAction(_ botsfw.WebhookContext, spaceID string) (m botsfw.MessageFromBot, err error) {
	if spaceID == "" {
		spaceID = "family"
	}

	var spaceIcon string

	var switchSpaceCallbackData string
	var switchSpaceTitle string
	switch spaceID {
	case "family":
		switchSpaceCallbackData = "space?id=private"
		switchSpaceTitle = "Private"
		spaceIcon = "👪"
	case "private":
		switchSpaceCallbackData = "space?id=family"
		switchSpaceTitle = "Family"
		spaceIcon = "🔒"
	}

	spaceTitle := strings.ToUpper(spaceID[:1]) + spaceID[1:]
	m.Text += fmt.Sprintf("Current space: %s <b>%s</b>", spaceIcon, spaceTitle)
	m.Format = botsfw.MessageFormatHTML

	firstRow := []tgbotapi.InlineKeyboardButton{
		{
			Text:         "📇 Contacts",
			CallbackData: "contacts",
		},
	}
	if spaceID != "private" {
		firstRow = append(firstRow, tgbotapi.InlineKeyboardButton{
			Text:         "👪 Members",
			CallbackData: "members",
		})
	}
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		firstRow,
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🚗 Assets",
				CallbackData: "assets",
			},
			{
				Text:         "💰 Budget",
				CallbackData: "budget",
			},
			{
				Text:         "💸 Debts",
				CallbackData: "debts",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🛒 Buy",
				CallbackData: "buy",
			},
			{
				Text:         "🏗️ ToDo",
				CallbackData: "todo",
			},
			{
				Text:         "📽️ Watch",
				CallbackData: "watch",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🗓️ Calendar",
				CallbackData: "calendar",
			},
			{
				Text:         "⚙️ Settings",
				CallbackData: "settings",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("🔀 Switch to \"%s\" space", switchSpaceTitle),
				CallbackData: switchSpaceCallbackData,
			},
		},
	)
	return
}
