package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/src/botscore/tghelpers"
	"net/url"
)

var assetsCommand = botsfw.Command{
	Code:     "assets",
	Commands: []string{"/assets"},
	InputTypes: []botsfw.WebhookInputType{
		botsfw.WebhookInputText,
		botsfw.WebhookInputCallbackQuery,
	},
	CallbackAction: assetsCallbackAction,
	Action:         assetsAction,
}

func assetsCallbackAction(whc botsfw.WebhookContext, callbackUrl *url.URL) (m botsfw.MessageFromBot, err error) {
	if m, err = assetsAction(whc); err != nil {
		return
	}

	keyboard := m.Keyboard.(*tgbotapi.InlineKeyboardMarkup)
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		tghelpers.BackToSpaceMenuButton(callbackUrl),
	})
	if m, err = whc.NewEditMessage(m.Text, m.Format); err != nil {
		return
	}
	m.Keyboard = keyboard

	m.EditMessageUID, err = tghelpers.GetEditMessageUID(whc)
	return
}

func assetsAction(_ botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
	m.Text = "<b>Family assets</b>"
	m.Text += "\n\n<i>Not implemented yet</i>"
	m.Format = botsfw.MessageFormatHTML
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			{
				Text: "💻 Manage in app",
				WebApp: &tgbotapi.WebappInfo{
					Url: "https://local-app.sneat.ws/space/family/h4qax/assets", // TODO: generate URL
				},
			},
			{
				Text:         "➕ Add asset",
				CallbackData: "/add-asset",
			},
		},
	)
	return
}
