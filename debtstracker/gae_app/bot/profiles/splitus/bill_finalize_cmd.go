package splitus

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/emoji"
	"github.com/sneat-co/sneat-go-backend/debtstracker/gae_app/debtstracker/models"
	"net/url"
)

const finalizeBillCommandCode = "finalize_bill"

var finalizeBillCommand = billCallbackCommand(finalizeBillCommandCode,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models.Bill) (m botsfw.MessageFromBot, err error) {
		footer := "<b>Are you ready to split the bill?</b>" +
			"\n" + "You won't be able to add/remove participants or change total once the bill is finalized."
		if m.Text, err = getBillCardMessageText(whc.Context(), whc.GetBotCode(), whc, bill, true, footer); err != nil {
			return
		}
		m.Format = botsfw.MessageFormatHTML
		m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         emoji.GREEN_CHECKBOX + " Yes, split the bill!",
					CallbackData: billCallbackCommandData(finalizeBillCommandCode, bill.ID),
				},
			},
			[]tgbotapi.InlineKeyboardButton{
				{
					Text:         emoji.NO_ENTRY_SIGN_ICON + " " + "Cancel",
					CallbackData: billCardCallbackCommandData(bill.ID),
				},
			},
		)
		return
	},
)
