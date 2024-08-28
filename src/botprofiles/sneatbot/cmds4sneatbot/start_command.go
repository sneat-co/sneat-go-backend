package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var startCommand = botsfw.Command{
	Code:       "start",
	Commands:   []string{"/start"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputText, botsfw.WebhookInputInlineQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
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
					CallbackData: "/members",
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
	},
}
