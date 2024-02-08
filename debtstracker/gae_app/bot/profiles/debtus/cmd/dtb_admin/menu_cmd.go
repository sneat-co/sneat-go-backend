package dtb_admin

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/bot/profiles/debtus/cmd/dtb_invite"
)

var AdminCommand = botsfw.Command{
	Code:     "admin",
	Commands: []string{"/admin"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m = whc.NewMessage("Admin menu")
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{Text: "Create mass invite", CallbackData: dtb_invite.CREATE_MASS_INVITE_CALLBACK_PATH},
			},
		)
		return
	},
}
