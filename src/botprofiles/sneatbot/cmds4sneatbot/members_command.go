package cmds4sneatbot

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/strongo/logus"
)

var membersCommand = botsfw.Command{
	Code:       "members",
	Commands:   []string{"/members"},
	InputTypes: []botsfw.WebhookInputType{botsfw.WebhookInputText, botsfw.WebhookInputCallbackQuery},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		ctx := whc.Context()
		logus.Infof(ctx, "membersCommand.Action(): InputType=%v", whc.Input().InputType())
		m.Text = "<b>Family members</b>"
		m.Format = botsfw.MessageFormatHTML
		//m.ResponseChannel = botsfw.BotAPISendMessageOverHTTPS // TODO: remove this line after debugging
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{
					Text: "Manage members in web app",
					//URL:  "https://local-app.sneat.ws/space/family",
					WebApp: &tgbotapi.WebappInfo{
						Url: "https://local-app.sneat.ws/space/family/h4qax/members", // TODO: generate URL
					},
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         "🧑 Myself",
					CallbackData: "/contact=myself",
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         "➕ Add family member",
					CallbackData: "/add-member",
				},
			},
		)
		m.ResponseChannel = botsfw.BotAPISendMessageOverHTTPS // TODO: remove this line after debugging
		return
	},
}
