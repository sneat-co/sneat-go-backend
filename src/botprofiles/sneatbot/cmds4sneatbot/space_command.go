package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"net/url"
)

var spaceCommand = botsfw.Command{
	Code:           "space",
	Commands:       []string{"/space"},
	InputTypes:     []botsfw.WebhookInputType{botsfw.WebhookInputCallbackQuery},
	CallbackAction: spaceCallbackAction,
}

func spaceCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	if m, err = spaceAction(whc); err != nil {
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

func spaceAction(_ botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	m.Text += "Current space: 👪 <b>Family</b>"
	m.Format = botsfw.MessageFormatHTML
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "📇 Contacts",
				CallbackData: "/contacts",
			},
			{
				Text:         "👪 Members",
				CallbackData: "members",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🚗 Assets",
				CallbackData: "/assets",
			},
			{
				Text:         "💰 Budget",
				CallbackData: "/budget",
			},
			{
				Text:         "💸 Debts",
				CallbackData: "/debtus",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🛒 Buy",
				CallbackData: "buy",
			},
			{
				Text:         "🏗️ ToDo",
				CallbackData: "do",
			},
			{
				Text:         "📽️ Watch",
				CallbackData: "watch",
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "🗓️ Calendar",
				CallbackData: "/calendar",
			},
			{
				Text:         "⚙️ Settings",
				CallbackData: "/settings",
			},
		},
	)
	return
}
