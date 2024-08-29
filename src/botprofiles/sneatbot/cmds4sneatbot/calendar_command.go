package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"net/url"
)

var calendarCommand = botsfw.Command{
	Code:     "calendar",
	Commands: []string{"/calendar"},
	InputTypes: []botsfw.WebhookInputType{
		botsfw.WebhookInputText,
		botsfw.WebhookInputCallbackQuery,
	},
	CallbackAction: calendarCallbackAction,
	Action:         calendarAction,
}

func calendarCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	if m, err = calendarAction(whc); err != nil {
		return
	}

	keyboard := m.Keyboard.(*tgbotapi.InlineKeyboardMarkup)
	spaceRef := tghelpers.GetSpaceRef(callbackUrl)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		tghelpers.BackToSpaceMenuButton(spaceRef),
	})
	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = keyboard

	m.EditMessageUID, err = tghelpers.GetEditMessageUID(whc)
	return
}

func calendarAction(_ botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	m.Format = botsfw.MessageFormatHTML
	m.Text = "<b>Family calendar</b>"
	m.Text += "\n\n<i>Not implemented yet</i>"
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "💻 Manage in app",
				WebApp: &tgbotapi.WebappInfo{
					Url: "https://local-app.sneat.ws/space/family/h4qax/calendar", // TODO: generate URL
				},
			},
			{
				Text: "➕ Add event",
				WebApp: &tgbotapi.WebappInfo{
					Url: "https://local-app.sneat.ws/space/family/h4qax/calendar", // TODO: generate URL
				},
			},
		},
		[]tgbotapi.InlineKeyboardButton{
			{
				Text:         "📆 Yesterday",
				CallbackData: "calendar?action=yesterday",
			},
			{
				Text:         "📅 Today",
				CallbackData: "calendar?action=today",
			},
			{
				Text:         "🗓️ Tomorrow",
				CallbackData: "calendar?action=tomorrow",
			},
		},
	)
	return
}
