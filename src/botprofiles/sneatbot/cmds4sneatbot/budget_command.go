package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"net/url"
)

var budgetCommand = botsfw.Command{
	Code:     "budget",
	Commands: []string{"/budget"},
	InputTypes: []botsfw.WebhookInputType{
		botsfw.WebhookInputText,
		botsfw.WebhookInputCallbackQuery,
	},
	CallbackAction: budgetCallbackAction,
	Action:         budgetAction,
}

func budgetCallbackAction(whc botsfw.WebhookContext, _ *url.URL) (m botsfw.MessageFromBot, err error) {
	if m, err = budgetAction(whc); err != nil {
		return
	}

	keyboard := m.Keyboard.(*tgbotapi.InlineKeyboardMarkup)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		tghelpers.BackToSpaceMenuButton(),
	})
	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = keyboard

	m.EditMessageUID, err = tghelpers.GetEditMessageUID(whc)
	return
}

func budgetAction(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	m.Text = "<b>Family budget</b>"
	m.Format = botsfw.MessageFormatHTML
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "💻 Manage in app",
				WebApp: &tgbotapi.WebappInfo{
					Url: "https://local-app.sneat.ws/space/family/h4qax/budget", // TODO: generate URL
				},
			},
		},
	)
	m.ResponseChannel = botsfw.BotAPISendMessageOverHTTPS // TODO: remove this line after debugging
	return
}
