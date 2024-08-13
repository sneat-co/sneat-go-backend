package botcmds4splitus

import (
	"github.com/bots-go-framework/bots-api-telegram/tgbotapi"
	"github.com/bots-go-framework/bots-fw/botsfw"
	"github.com/dal-go/dalgo/dal"
	"github.com/sneat-co/debtstracker-translations/trans"
	"github.com/sneat-co/sneat-go-backend/src/modules/splitus/models4splitus"
	"github.com/strongo/logus"
	"net/url"
)

const CHANGE_BILL_PAYER_COMMAND = "change-bill-payer"

var changeBillPayerCommand = billCallbackCommand(CHANGE_BILL_PAYER_COMMAND,
	func(whc botsfw.WebhookContext, _ dal.ReadwriteTransaction, callbackUrl *url.URL, bill models4splitus.BillEntry) (m botsfw.MessageFromBot, err error) {
		c := whc.Context()
		logus.Debugf(c, "changeBillPayerCommand.CallbackAction()")
		var (
			mt string
			//editedMessage *tgbotapi.EditMessageTextConfig
		)
		if mt, err = getBillCardMessageText(c, whc.GetBotCode(), whc, bill, true, whc.Translate(trans.MESSAGE_TEXT_BILL_ASK_WHO_PAID)); err != nil {
			return
		}
		if m, err = whc.NewEditMessage(mt, botsfw.MessageFormatHTML); err != nil {
			return
		}
		markup := tgbotapi.NewInlineKeyboardMarkup()

		for _, member := range bill.Data.GetBillMembers() {
			s := member.Name
			if member.Paid > 0 {
				s = "✔ " + s
			}

			markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
				{
					Text:         s,
					CallbackData: billCardCallbackCommandData(bill.ID),
				},
			})
		}

		markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
			{
				Text:         whc.Translate(trans.BUTTON_TEXT_CANCEL),
				CallbackData: billCardCallbackCommandData(bill.ID),
			},
		})

		m.Keyboard = markup
		return
	},
)
