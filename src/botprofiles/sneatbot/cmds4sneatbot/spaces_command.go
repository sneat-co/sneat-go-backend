package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
)

var spacesCommand = botsfw.Command{
	Code:       "spaces",
	Commands:   []string{"/spaces"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputText},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Text = "<b>Your spaces</b>"
		m.Text += "\nCurrent space: <b>Family</b>"
		m.Text += "\nClick to switch current space."

		m.Format = botsfw.MessageFormatHTML
		//m.Keyboard = tgbotapi.NewReplyKeyboard(
		//	[]tgbotapi.KeyboardButton{
		//		{
		//			Text: "👪 Family",
		//		},
		//	},
		//	[]tgbotapi.KeyboardButton{
		//		{
		//			Text: "➕ Add new space",
		//		},
		//	},
		//)
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         "👪 Family ✅",
					CallbackData: "/space=family",
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         "🔒 Private 🔲",
					CallbackData: "/space=family",
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         "➕ Add new space (not implemented yet)",
					CallbackData: "/add-space",
				},
			},
		)
		return
	},
}
