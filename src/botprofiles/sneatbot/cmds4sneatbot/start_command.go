package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"net/url"
)

var startCommand = botsfw.Command{
	Code:     "start",
	Commands: []string{"/start"},
	InputTypes: []botsfw.WebhookInputType{
		botsfw.WebhookInputText,
		botsfw.WebhookInputCallbackQuery,
		botsfw.WebhookInputInlineQuery,
	},
	Action: startAction,
	CallbackAction: func(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
		if m, err = startAction(whc); err != nil {
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
	},
}

func startAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	m.Text = "Hello, stranger!"
	m.Text += "\n\nI'm a @SneatBot. I can help you to manage your day-to-day family life."
	m.Text += "\n\nOr you can create a space to manage your group/team/community."
	m.Text += "\n\nCurrent space: 👪 <b>Family</b>"
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
				Text:         "🛒 To-Buy",
				CallbackData: "/listus to-buy",
			},
			{
				Text:         "🏗️ To-Do",
				CallbackData: "/listus to-do",
			},
			{
				Text:         "📽️ To-Watch",
				CallbackData: "/listus to-watch",
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
