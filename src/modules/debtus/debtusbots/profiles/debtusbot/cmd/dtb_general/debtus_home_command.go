package dtb_general

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/debtstracker-translations/trans"
)

const DebtusHomeCommandCode = "debtus"

var DebtusHomeCommand = botsfw.Command{
	Code:     DebtusHomeCommandCode,
	Commands: []string{"/debtus"},
	Action: func(whc botsfw.WebhookContext) (m botsfw.MessageFromBot, err error) {
		m.Format = botsfw.MessageFormatHTML
		m.Text = `<b>Debtus home</b>

Choose what you want to do:
`
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData(
					whc.CommandText(trans.COMMAND_TEXT_GAVE, emoji.GIVE_ICON),
					"balance",
				),
				tgbotapi.NewInlineKeyboardButtonData(
					whc.CommandText(trans.COMMAND_TEXT_GOT, emoji.TAKE_ICON),
					"balance",
				),
			},
			[]tgbotapi.InlineKeyboardButton{
				tgbotapi.NewInlineKeyboardButtonData(
					whc.CommandText(trans.COMMAND_TEXT_BALANCE, emoji.BALANCE_ICON),
					"balance",
				),
				tgbotapi.NewInlineKeyboardButtonData(
					whc.CommandText(trans.COMMAND_TEXT_HISTORY, emoji.HISTORY_ICON),
					"balance",
				),
			},
		)
		return
	},
}
